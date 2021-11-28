package executioner

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Imputes/fdlr/internal/downloader"
	"github.com/Imputes/fdlr/internal/merger"
	"github.com/Imputes/fdlr/internal/tool"

	"github.com/pkg/errors"
)

func Do(url string, state *tool.State, conc int) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	files := []string{}
	parts := []tool.DownloadRange{}
	isInterrupted := false

	doneChan := make(chan bool, conc)
	fileChan := make(chan string, conc)
	errorChan := make(chan error, 1)
	stateChan := make(chan tool.DownloadRange, 1)
	interruptChan := make(chan bool, conc)

	var dlr *downloader.HTTPDownloader

	var err error

	if state == nil {
		dlr, err = downloader.NewHTTPDownloader(url, conc)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		dlr = &downloader.HTTPDownloader{
			URL:            state.URL,
			File:           filepath.Base(state.URL),
			Part:           int64(len(state.DownloadRanges)),
			SkipTLS:        true,
			DownloadRanges: state.DownloadRanges,
			Resumable:      true,
		}
	}

	// start downloading
	go dlr.Downloading(doneChan, fileChan, errorChan, interruptChan, stateChan)

	// monitoring the download in progress
	for {
		select {
		case <-signalChan:
			isInterrupted = true
			for conc > 0 {
				interruptChan <- true
				conc--
			}
		case file := <-fileChan:
			files = append(files, file)
		case err = <-errorChan:
			return errors.WithStack(err)
		case part := <-stateChan:
			parts = append(parts, part)
		case <-doneChan:
			if isInterrupted {
				if dlr.Resumable {
					fmt.Printf("Interrupted, saving state ... \n")
					s := &tool.State{
						URL:            url,
						DownloadRanges: parts,
					}
					if err = s.Save(); err != nil {
						return errors.WithStack(err)
					}
					return nil
				} else {
					fmt.Printf("Interrupted, but downloading url is not resumable, silently die\n")
					return nil
				}
			} else {
				err = merger.MergeFile(files, filepath.Base(url))
				if err != nil {
					return errors.WithStack(err)
				}

				folder, err := tool.GetFolderFrom(url)
				if err != nil {
					return errors.WithStack(err)
				}

				err = os.RemoveAll(folder)
				if err != nil {
					return errors.WithStack(err)
				}

				return nil
			}
		}
	}
}
