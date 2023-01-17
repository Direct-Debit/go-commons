package httpext

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	address := "0.0.0.0:12345"
	s := RunTestServer(address)
	fmt.Println("Giving test server some time to spin up")
	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		go func() {
			resp, err := Get("http://" + address + "/wait")

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()

		openFiles := func() int {
			out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
			if err != nil {
				return 0
			}
			lines := strings.Split(string(out), "\n")
			return len(lines) - 1
		}()
		fmt.Println(openFiles)
	}

	time.Sleep(10 * time.Second)
	_ = s.Shutdown(context.Background())

}
