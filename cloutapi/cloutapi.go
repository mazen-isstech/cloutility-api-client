package cloutapi

/* Copyright 2022-2023 (C) Blue Safespring AB
   Programmed by Jan Johansson
   Contributions by Daniel de Oquiñena and Patrik Lundin
   All rights reserved for now, will have liberal
   license later */

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AuthenticatedClient struct {
	HttpClient   *http.Client
	BaseURL      string
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Expires      int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// Temporary function for testing purposes
func RunClient() {
	var (
		user me
		// myConsumer consumer
		// myNode     node
	)

	// Initialize client by passing username, password and client_id from config file
	c, err := Init(
		context.TODO(),
		viper.GetString("client_id"),
		viper.GetString("client_origin"),
		viper.GetString("username"),
		viper.GetString("password"),
		viper.GetString("url"),
	)
	if err != nil {
		log.Printf("Error authenticating: %s", err)
		os.Exit(1)
	}

	// if viper.GetBool("debug") {
	// 	log.Println("Token type:", c.TokenType)
	// 	log.Println("Expires:", c.Expires)
	// 	log.Println("Refresh token:", c.RefreshToken)
	// 	log.Println("Access token:", c.AccessToken)
	// }

	user, err = c.GetUser()
	if err != nil {
		log.Println("Error retrieving userdata: ", err)
	}

	fmt.Println("USER: ", user, "\n\n")

	// if viper.GetBool("debug") {
	// 	log.Println(user.Name)
	// 	log.Println(user.BusinessUnit.Name)
	// 	log.Println(user.BusinessUnit.ID)
	// }

	// node, err := c.GetNode(user.BusinessUnitID, user.ID)
	// if err != nil {
	// 	log.Println("Error retrieving nodedata: ", err)
	// }n
	// fmt.Println("NODE1: ", node, "\n\n")

	if viper.GetBool("dry-run") {
		log.Println("Running in dry-run mode, exiting")
		os.Exit(0)
	}

	consumer, err := c.CreateConsumer(user.BusinessUnit.ID, "testar")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(user.BusinessUnit.ID, " ", consumer.ID)
	fmt.Println("CONSUMER: ", consumer, "\n\n")
	node, err := c.CreateNode(user.BusinessUnit.ID, consumer.ID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("NODE2: ", node, "\n\n")
	node, err = c.DeleteNode(node.ID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("NODE DELETED: ", node, "\n\n")
	// log.Println(myNode)
	fmt.Println(c.AccessToken)
	ok, err := c.DeleteConsumer(user.BusinessUnit.ID, consumer.ID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(user.BusinessUnit.ID, " ", consumer.ID)
	fmt.Println("DETETED: ", ok, "\n\n")
}

func (c *AuthenticatedClient) apiRequest(endpoint string, method string, payload []byte) (string, error) {
	ctx := context.Background()

	var reader io.Reader
	if payload != nil {
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return "", fmt.Errorf("failed to complete request: %s", err)
	}

	req.Header.Set("User-Agent", "safespring-golang-client")
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Origin", viper.GetString("client_origin"))
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	// XXX - needs conf file

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve response body: %s", err)
	}
	defer resp.Body.Close()

	// Check response code and return error if not 2xx
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return "", fmt.Errorf("error response %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	return string(body), nil
}

// Initialize client and return an AuthenticatedClient
func Init(ctx context.Context, client_id, origin, username, password, baseURL string) (*AuthenticatedClient, error) {
	var c AuthenticatedClient

	// Initialize http.Client
	c.HttpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// Get baseurl from passed variable
	c.BaseURL = baseURL

	authurl := "/v1/oauth"

	// Construct body
	loginData := url.Values{}
	loginData.Add("client_id", client_id)
	loginData.Add("grant_type", "password")
	loginData.Add("username", username)
	loginData.Add("password", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+authurl, strings.NewReader(loginData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set header propertys
	req.Header.Set("User-Agent", "safespring-golang-client")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", origin)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete authentication request: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to retrieve authentication: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read authentication response: %s", err)
	}

	if err := json.Unmarshal([]byte(body), &c); err != nil {
		return nil, fmt.Errorf("failed to decode authentication response: %s", err)
	}

	return &c, nil
}
