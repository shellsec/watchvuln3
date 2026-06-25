package ctrl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoardPublicURL(t *testing.T) {
	assert := require.New(t)

	cases := []struct {
		name   string
		config WatchVulnAppConfig
		want   string
	}{
		{
			name:   "no board",
			config: WatchVulnAppConfig{},
			want:   "",
		},
		{
			name: "full public url override",
			config: WatchVulnAppConfig{
				WebAddr:      "0.0.0.0:8765",
				WebPublicURL: "http://vuln.example.com:9000/board",
			},
			want: "http://vuln.example.com:9000/board",
		},
		{
			name: "public host with port from web_addr",
			config: WatchVulnAppConfig{
				WebAddr:       "0.0.0.0:8765",
				WebPublicHost: "192.168.1.100",
			},
			want: "http://192.168.1.100:8765/",
		},
		{
			name: "port change follows web_addr",
			config: WatchVulnAppConfig{
				WebAddr:       "0.0.0.0:9000",
				WebPublicHost: "192.168.1.100",
			},
			want: "http://192.168.1.100:9000/",
		},
		{
			name: "loopback without public host",
			config: WatchVulnAppConfig{
				WebAddr: "127.0.0.1:8765",
			},
			want: "http://127.0.0.1:8765/",
		},
		{
			name: "wildcard listen without public host",
			config: WatchVulnAppConfig{
				WebAddr: "0.0.0.0:8765",
			},
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(tc.want, tc.config.BoardPublicURL())
		})
	}
}

func TestBoardFeedPublicURL(t *testing.T) {
	assert := require.New(t)

	cases := []struct {
		name   string
		config WatchVulnAppConfig
		want   string
	}{
		{
			name: "from public host and web_addr port",
			config: WatchVulnAppConfig{
				WebAddr:       "0.0.0.0:8766",
				WebPublicHost: "192.168.1.100",
			},
			want: "http://192.168.1.100:8766/feed.xml",
		},
		{
			name: "from full public url",
			config: WatchVulnAppConfig{
				WebAddr:      "0.0.0.0:8765",
				WebPublicURL: "http://vuln.example.com:9000/",
			},
			want: "http://vuln.example.com:9000/feed.xml",
		},
		{
			name: "loopback",
			config: WatchVulnAppConfig{
				WebAddr: "127.0.0.1:8765",
			},
			want: "http://127.0.0.1:8765/feed.xml",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(tc.want, tc.config.BoardFeedPublicURL())
		})
	}
}
