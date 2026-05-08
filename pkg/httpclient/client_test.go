package httpclient

// func TestClient(t *testing.T) {

// 	var params []any
// 	jsonReq := map[string]any{
// 		"jsonrpc": "3.0",
// 		"method":  "foo.bar",
// 		"params":  params,
// 		"id":      12,
// 	}
// 	data, _ := json.Marshal(jsonReq)

// 	setterFunc := func(req *http.Request) {
// 		req.Header.Set("Content-Length", strconv.Itoa(len(data)))
// 		req.Header.Set("Content-Type", "application/json")
// 		// req.Header.Set("Content-Type", "multipart/form-data;charset=utf8;")
// 		req.Header.Set("Accept", "application/json")
// 	}

// 	myHttpClient := NewHttpClient()
// 	_, _, statusCode, err := myHttpClient.Req("https://foo.example.com/api/v1/example", "POST", data, setterFunc)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(statusCode)
// 	}
// }
