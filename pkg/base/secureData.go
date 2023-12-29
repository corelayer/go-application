/*
 * Copyright 2023 CoreLayer BV
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package base

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/minio/sio"
	"golang.org/x/crypto/hkdf"
)

type SecureData struct {
	Nonce       string `json:"nonce" yaml:"nonce" mapstructure:"nonce"`
	CipherSuite string `json:"ciphersuite" yaml:"ciphersuite" mapstructure:"ciphersuite"`
	HexData     string `json:"hexdata" yaml:"hexdata" mapstructure:"hexdata"`
}

func (d *SecureData) Encrypt(master string) error {
	var (
		err          error
		source       []byte
		cryptoConfig sio.Config
	)

	cryptoConfig, err = d.getCryptoConfig(master)
	if err != nil {
		return fmt.Errorf("could not initialize crypto parameters: %w", err)
	}

	source, err = d.Bytes()
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	outBuf := make([]byte, 0)
	encryptedData := bytes.NewBuffer(outBuf)
	if _, err = sio.Encrypt(encryptedData, bytes.NewReader(source), cryptoConfig); err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	d.HexData = hex.EncodeToString(encryptedData.Bytes())
	return nil
}

func (d *SecureData) Decrypt(master string) error {
	var (
		err          error
		source       []byte
		cryptoConfig sio.Config
	)

	if d.Nonce == "" {
		return fmt.Errorf("nonce is not set, cannot decrypt")
	}

	cryptoConfig, err = d.getCryptoConfig(master)
	if err != nil {
		return fmt.Errorf("could not initialize crypto parameters: %w", err)
	}

	source, err = d.Bytes()
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	outBuf := make([]byte, 0)
	decryptedData := bytes.NewBuffer(outBuf)
	if _, err = sio.Decrypt(decryptedData, bytes.NewReader(source), cryptoConfig); err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Reset nonce after decryption, to make sure it gets rotated on next encryption
	d.resetNonce()
	d.HexData = hex.EncodeToString(decryptedData.Bytes())
	return nil
}

func (d *SecureData) Update(data []byte) error {
	if d.Nonce != "" {
		return fmt.Errorf("cannot update encrypted data")
	}

	d.HexData = hex.EncodeToString(data)
	return nil
}

func (d *SecureData) Bytes() ([]byte, error) {
	return hex.DecodeString(d.HexData)
}

func (d *SecureData) getCryptoConfig(master string) (sio.Config, error) {
	var (
		err          error
		masterKey    []byte
		nonce        []byte
		key          [32]byte
		cipherSuites []byte
	)

	masterKey, err = hex.DecodeString(master)
	if err != nil {
		return sio.Config{}, fmt.Errorf("could not decode master key: %w", err)
	}

	nonce, err = d.getNonce()
	if err != nil {
		return sio.Config{}, err
	}

	kdf := hkdf.New(sha256.New, masterKey, nonce, nil)
	if _, err = io.ReadFull(kdf, key[:]); err != nil {
		return sio.Config{}, fmt.Errorf("failed to derive encryption key: %w", err)
	}

	cipherSuites, err = d.getCipherSuite()
	if err != nil {
		return sio.Config{}, err
	}

	return sio.Config{Key: key[:], CipherSuites: cipherSuites}, nil
}

func (d *SecureData) getNonce() ([]byte, error) {
	if d.Nonce != "" {
		return hex.DecodeString(d.Nonce)
	}

	var nonce [32]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to read random data for nonce: %w", err)
	}
	d.Nonce = hex.EncodeToString(nonce[:])
	return nonce[:], nil
}

func (d *SecureData) resetNonce() {
	d.Nonce = ""
}

func (d *SecureData) getCipherSuite() ([]byte, error) {
	switch d.CipherSuite {
	case "AES_256_GCM":
		return []byte{sio.AES_256_GCM}, nil
	case "CHACHA20_POLY1305":
		return []byte{sio.CHACHA20_POLY1305}, nil
	default:
		return nil, fmt.Errorf("invalid cipher suite %s", d.CipherSuite)
	}
}
