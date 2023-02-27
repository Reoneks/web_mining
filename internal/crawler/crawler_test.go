package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestCrawler_ParsePage(t *testing.T) {
	type args struct {
		page    io.Reader
		baseURL *url.URL
	}

	site := "https://stackoverflow.com/questions/40643030/how-to-get-webpage-content-into-a-string-using-go"
	resp, err := resty.New().SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/110.0").R().Get(site)
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse(site)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		c       *Crawler
		args    args
		wantErr bool
	}{
		{
			name:    "Test",
			args:    args{bytes.NewReader(resp.Body()), u},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{}
			got, err := c.ParsePage(tt.args.page, tt.args.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawler.ParsePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Println(got)
		})
	}
}
