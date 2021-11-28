package downloader

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	stdurl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Imputes/fdlr/internal/tool"

	"github.com/cheggaaa/pb"
	"github.com/pkg/errors"
)

const (
	acceptRange   = "Accept-Ranges"
	contentLength = "Content-Length"
)

var (
	client = &http.Client{Transport: tr}

	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	skipTLS = true
)

type HTTPDownloader struct {
	URL            string
	File           string
	Part           int64
	Len            int64
	IPs            []string
	SkipTLS        bool
	DownloadRanges []tool.DownloadRange
	Resumable      bool
}

func NewHTTPDownloader(url string, par int) (*HTTPDownloader, error) {
	resumable := true

	parsed, err := stdurl.Parse(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ips, err := net.LookupIP(parsed.Host)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ipstr := tool.GetIPv4(ips)
	fmt.Printf("Downloading IP is: %s\n", strings.Join(ipstr, " | "))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.Header.Get(acceptRange) == "" {
		fmt.Printf("This file does not support HTTP Range, so set c to 1\n")
		par = 1
	}

	// get download range
	clen := resp.Header.Get(contentLength)
	if clen == "" {
		fmt.Printf("Target url not contain Content-Length header, fallback to parallel 1\n")
		clen = "1" // cheggaaa/pb not accept 0 length
		par = 1
		resumable = false
	}

	fmt.Printf("Start downloading with %d connections \n", par)

	len, err := strconv.ParseInt(clen, 10, 64)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sizeInMb := float64(len) / (1024 * 1024)
	if clen == "1" {
		fmt.Printf("Download size: not specified\n")
	} else if sizeInMb < 1024 {
		fmt.Printf("Download target size: %.1f MB\n", sizeInMb)
	} else {
		fmt.Printf("Download target size: %.1f GB\n", sizeInMb/1024)
	}

	file := filepath.Base(url)
	ret := new(HTTPDownloader)
	ret.URL = url
	ret.File = file
	ret.Part = int64(par)
	ret.Len = len
	ret.IPs = ipstr
	ret.SkipTLS = skipTLS
	ret.DownloadRanges, err = partCalculate(int64(par), len, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret.Resumable = resumable

	return ret, nil
}

func partCalculate(par int64, len int64, url string) ([]tool.DownloadRange, error) {
	ret := []tool.DownloadRange{}

	for i := int64(0); i < par; i++ {
		from := (len / par) * i
		to := int64(0)
		if i < par-1 {
			to = (len/par)*(i+1) - 1
		} else {
			to = len
		}

		file := filepath.Base(url)
		folder, err := tool.GetFolderFrom(url)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := tool.Mkdir(folder); err != nil {
			return nil, errors.WithStack(err)
		}

		fname := fmt.Sprintf("%s.part%d", file, i)
		path := filepath.Join(folder, fname)
		ret = append(ret, tool.DownloadRange{
			URL:       url,
			Path:      path,
			RangeFrom: from,
			RangeTo:   to,
		})
	}
	return ret, nil
}

func (d *HTTPDownloader) Downloading(doneChan chan bool, fileChan chan string, errorChan chan error, interruptChan chan bool, stateSaveChan chan tool.DownloadRange) {
	var bars []*pb.ProgressBar
	var barpool *pb.Pool
	var err error

	if tool.DisappearProgressBar() {
		bars = []*pb.ProgressBar{}
		for i, part := range d.DownloadRanges {
			newbar := pb.New64(part.RangeTo - part.RangeFrom).SetUnits(pb.U_BYTES).Prefix(fmt.Sprintf("%s - %d", d.File, i))
			newbar.ShowBar = false
			bars = append(bars, newbar)
		}
		barpool, err = pb.StartPool(bars...)
		if err != nil {
			errorChan <- errors.WithStack(err)
			return
		}
	}

	// Parallel download
	ws := new(sync.WaitGroup)
	for i, p := range d.DownloadRanges {
		ws.Add(1)
		go func(d *HTTPDownloader, loop int64, part tool.DownloadRange) {
			defer ws.Done()
			bar := new(pb.ProgressBar)

			if tool.DisappearProgressBar() {
				bar = bars[loop]
			}

			ranges := ""
			if part.RangeTo != d.Len {
				ranges = fmt.Sprintf("bytes=%d-%d", part.RangeFrom, part.RangeTo)
			} else {
				ranges = fmt.Sprintf("bytes=%d-", part.RangeFrom)
			}

			req, err := http.NewRequest("GET", d.URL, nil)
			if err != nil {
				errorChan <- errors.WithStack(err)
				return
			}

			if d.Part > 1 {
				req.Header.Add("Range", ranges)
				if err != nil {
					errorChan <- errors.WithStack(err)
					return
				}
			}

			resp, err := client.Do(req)
			if err != nil {
				errorChan <- errors.WithStack(err)
				return
			}
			defer resp.Body.Close()

			f, err := os.OpenFile(part.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if err != nil {
				errorChan <- errors.WithStack(err)
				return
			}

			var writer io.Writer
			if tool.DisappearProgressBar() {
				writer = io.MultiWriter(f, bar)
			} else {
				writer = io.MultiWriter(f)
			}

			// copy 100 bytes each loop
			current := int64(0)
			for {
				select {
				case <-interruptChan:
					stateSaveChan <- tool.DownloadRange{
						URL:       d.URL,
						Path:      part.Path,
						RangeFrom: current + part.RangeFrom,
						RangeTo:   part.RangeTo,
					}
					return
				default:
					written, err := io.CopyN(writer, resp.Body, 100)
					current += written
					if err != nil {
						if err != io.EOF {
							errorChan <- errors.WithStack(err)
							return
						}
						fileChan <- part.Path
						return
					}
				}
			}
			err = f.Close()
			if err != nil {
				errorChan <- errors.WithStack(err)
				return
			}
		}(d, int64(i), p)
	}

	ws.Wait()

	err = barpool.Stop()
	if err != nil {
		errorChan <- errors.WithStack(err)
		return
	}

	doneChan <- true
}
