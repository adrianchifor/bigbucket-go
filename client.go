package bigbucket

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Client struct {
	address string
	timeout int
	gcpAuth bool
	headers map[string]string
	jwt     string
}

type ClientOption func(c *Client)

// Takes WithAddress(address string), WithTimeout(timeout int), WithGcpAuth(enabled bool)
// and WithRequestHeaders(headers map[string]string) as optional arguments
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		address: "http://localhost:8080",
		timeout: 30,
		gcpAuth: false,
		headers: make(map[string]string),
		jwt:     "",
	}
	for _, opt := range opts {
		opt(client)
	}

	return client
}

func WithAddress(address string) ClientOption {
	return func(c *Client) {
		c.address = address
	}
}

func WithTimeout(timeout int) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithGcpAuth(enabled bool) ClientOption {
	return func(c *Client) {
		c.gcpAuth = enabled
	}
}

func WithRequestHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

func (c *Client) GetTables() ([]string, error) {
	resp, err := httpRequest("GET", c, "/api/table", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, constructError(resp)
	}

	var data map[string][]string
	json.NewDecoder(resp.Body).Decode(&data)

	return data["tables"], nil
}

type Table struct {
	name   string
	client *Client
}

func (c *Client) UseTable(table string) *Table {
	return &Table{
		name:   table,
		client: c,
	}
}

func (t *Table) DeleteTable() error {
	params := map[string]string{
		"table": t.name,
	}
	resp, err := httpRequest("DELETE", t.client, "/api/table", params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return constructError(resp)
	}

	return nil
}

func (t *Table) ListColumns() ([]string, error) {
	params := map[string]string{
		"table": t.name,
	}
	resp, err := httpRequest("GET", t.client, "/api/column", params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, constructError(resp)
	}

	var data map[string][]string
	json.NewDecoder(resp.Body).Decode(&data)

	return data["columns"], nil
}

func (t *Table) DeleteColumn(column string) error {
	params := map[string]string{
		"table":  t.name,
		"column": column,
	}
	resp, err := httpRequest("DELETE", t.client, "/api/column", params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return constructError(resp)
	}

	return nil
}

type RequestConfig struct {
	prefix  string
	columns []string
	limit   string
}

type RequestConfigOption func(c *RequestConfig)

func WithPrefix(prefix string) RequestConfigOption {
	return func(rc *RequestConfig) {
		rc.prefix = prefix
	}
}

func WithColumns(columns []string) RequestConfigOption {
	return func(rc *RequestConfig) {
		rc.columns = columns
	}
}

func WithLimit(limit int) RequestConfigOption {
	return func(rc *RequestConfig) {
		rc.limit = strconv.Itoa(limit)
	}
}

func populateRequestConfig(opts []RequestConfigOption) *RequestConfig {
	rc := &RequestConfig{
		prefix:  "",
		columns: []string{},
		limit:   "",
	}
	for _, opt := range opts {
		opt(rc)
	}

	return rc
}

// Takes WithPrefix(prefix string) as an optional argument
func (t *Table) CountRows(opts ...RequestConfigOption) (int64, error) {
	rc := populateRequestConfig(opts)
	params := map[string]string{
		"table": t.name,
	}
	if rc.prefix != "" {
		params["prefix"] = rc.prefix
	}
	resp, err := httpRequest("GET", t.client, "/api/row/count", params, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, constructError(resp)
	}

	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	rowsCount, _ := strconv.ParseInt(data["rowsCount"], 0, 64)

	return rowsCount, nil
}

// Takes WithPrefix(prefix string) as an optional argument
func (t *Table) ListRows(opts ...RequestConfigOption) ([]string, error) {
	rc := populateRequestConfig(opts)
	params := map[string]string{
		"table": t.name,
	}
	if rc.prefix != "" {
		params["prefix"] = rc.prefix
	}
	resp, err := httpRequest("GET", t.client, "/api/row/list", params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, constructError(resp)
	}

	var data map[string][]string
	json.NewDecoder(resp.Body).Decode(&data)

	return data["rowKeys"], nil
}

// Takes WithColumns(columns []string) as an optional argument
func (t *Table) ReadRow(key string, opts ...RequestConfigOption) (map[string]string, error) {
	rc := populateRequestConfig(opts)
	params := map[string]string{
		"table": t.name,
		"key":   key,
	}
	if len(rc.columns) > 0 {
		params["columns"] = strings.Join(rc.columns, ",")
	}
	resp, err := httpRequest("GET", t.client, "/api/row", params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, constructError(resp)
	}

	var data map[string]map[string]string
	json.NewDecoder(resp.Body).Decode(&data)

	return data[key], nil
}

// Takes WithPrefix(prefix string), WithColumns(columns []string) and WithLimit(limit int) as optional arguments
func (t *Table) ReadRows(opts ...RequestConfigOption) (map[string]map[string]string, error) {
	rc := populateRequestConfig(opts)
	params := map[string]string{
		"table": t.name,
	}
	if rc.prefix != "" {
		params["prefix"] = rc.prefix
	}
	if len(rc.columns) > 0 {
		params["columns"] = strings.Join(rc.columns, ",")
	}
	if rc.limit != "" {
		params["limit"] = rc.limit
	}
	resp, err := httpRequest("GET", t.client, "/api/row", params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, constructError(resp)
	}

	var data map[string]map[string]string
	json.NewDecoder(resp.Body).Decode(&data)

	return data, nil
}

func (t *Table) SetRow(key string, columnValueMap map[string]string) error {
	params := map[string]string{
		"table": t.name,
		"key":   key,
	}
	resp, err := httpRequest("POST", t.client, "/api/row", params, columnValueMap)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return constructError(resp)
	}

	return nil
}

func (t *Table) DeleteRow(key string) error {
	params := map[string]string{
		"table": t.name,
		"key":   key,
	}
	resp, err := httpRequest("DELETE", t.client, "/api/row", params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return constructError(resp)
	}

	return nil
}

func (t *Table) DeleteRows(prefix string) error {
	params := map[string]string{
		"table":  t.name,
		"prefix": prefix,
	}
	resp, err := httpRequest("DELETE", t.client, "/api/row", params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return constructError(resp)
	}

	return nil
}
