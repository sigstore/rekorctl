// Copyright 2021 The Sigstore Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	tclient "github.com/google/trillian/client"
	tcrypto "github.com/google/trillian/crypto"
	"github.com/google/trillian/merkle/logverifier"
	"github.com/google/trillian/merkle/rfc6962/hasher"

	"github.com/sigstore/rekor/pkg/log"
)

func DoGet(url string, rekorEntry []byte) error {
	log := log.Logger
	// Set Context with Timeout for connects to thde log rpc server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := addFileToRequest(request, bytes.NewReader(rekorEntry)); err != nil {
		log.Error(err)
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	resp := getProofResponse{}
	if err := json.Unmarshal(content, &resp); err != nil {
		log.Error(err)
		return err
	}

	pub, err := x509.ParsePKIXPublicKey(resp.Key)
	if err != nil {
		log.Error(err)
		return err
	}

	if resp.Proof != nil {
		leafHash := hasher.DefaultHasher.HashLeaf(rekorEntry)
		verifier := tclient.NewLogVerifier(hasher.DefaultHasher, pub, crypto.SHA256)
		root, err := tcrypto.VerifySignedLogRoot(verifier.PubKey, verifier.SigHash, resp.Proof.SignedLogRoot)
		if err != nil {
			log.Error(err)
			return err
		}

		v := logverifier.New(hasher.DefaultHasher)
		proof := resp.Proof.Proof[0]
		if err := v.VerifyInclusionProof(proof.LeafIndex, int64(root.TreeSize), proof.Hashes, root.RootHash, leafHash); err != nil {
			log.Error(err)
			return err
		}

		log.Info("Proof correct!")
	} else {
		log.Info(resp.Status)
	}

	return nil
}

func addFileToRequest(request *http.Request, r io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, err := writer.CreateFormFile("fileupload", "linkfile")
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, r); err != nil {
		return err
	}

	request.Body = ioutil.NopCloser(body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	return nil
}
