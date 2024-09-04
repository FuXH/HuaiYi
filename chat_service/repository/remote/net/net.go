package net

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HttpClientGet(ctx context.Context,
	urlPath string, param map[string]interface{}, header map[string]string,
	rsp interface{}) error {
	client := &http.Client{}

	urlPath = encodeParam(urlPath, param)
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	fmt.Println(string(body))
	if err := json.Unmarshal(body, rsp); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func HttpClientPost(ctx context.Context,
	urlPath string, param map[string]interface{}, header map[string]string,
	req, rsp interface{}) error {
	client := &http.Client{}

	urlPath = encodeParam(urlPath, param)
	reqBody, _ := json.Marshal(req)
	reqCli, err := http.NewRequest("POST", urlPath, strings.NewReader(string(reqBody)))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	for k, v := range header {
		reqCli.Header.Add(k, v)
	}

	resp, err := client.Do(reqCli)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if err := json.Unmarshal(body, rsp); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func encodeParam(urlPath string, param map[string]interface{}) string {
	reqParam := url.Values{}
	for k, v := range param {
		str := fmt.Sprintf("%v", v)
		reqParam.Add(k, str)
	}
	if len(reqParam.Encode()) != 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, reqParam.Encode())
	}
	return urlPath
}
