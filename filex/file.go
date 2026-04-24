package filex

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/soyacen/gox/stringx"
)

// Lines splits a file into lines and returns an iterator sequence, yielding one line per iteration.
//
// Parameters:
//   - file: The file object to read
//
// Returns:
//   - iter.Seq[[]byte]: An iterator sequence where each element is a line as a byte slice
func Lines(file *os.File) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// 检查扫描过程中是否出现错误
			if scanner.Err() != nil {
				return
			}
			// 将当前行数据传递给yield函数，并检查是否需要继续迭代
			if !yield(scanner.Bytes()) {
				return
			}
		}
	}
}

// Primary returns the primary name of a file without its extension.
//
// Parameters:
//   - path: The file path
//
// Returns:
//   - string: The file name without extension
func Primary(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// List lists the paths of all regular files under the specified root directory.
//
// Parameters:
//   - root: The root directory to walk through
//
// Returns:
//   - []string: List of file paths
//   - error: Error encountered during directory walk
func List(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

// Extension returns the file extension without the leading dot.
//
// Parameters:
//   - path: The file path
//
// Returns:
//   - string: The file extension
func Extension(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}

var DownloadClient = &http.Client{}

// DownloadToReader downloads content from a URL and returns it as a ReadCloser.
//
// Parameters:
//   - ctx: The context for cancellation
//   - url: The URL to download from
//
// Returns:
//   - io.ReadCloser: The response body
//   - error: Error if the download fails
func DownloadToReader(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := DownloadClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return resp.Body, nil
}

// DownloadToData downloads content from a URL and returns it as a byte slice.
//
// Parameters:
//   - ctx: The context for cancellation
//   - url: The URL to download from
//
// Returns:
//   - []byte: The downloaded content
//   - error: Error if the download fails
func DownloadToData(ctx context.Context, url string) ([]byte, error) {
	reader, err := DownloadToReader(ctx, url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()
	return io.ReadAll(reader)
}

// DownloadToFile downloads content from a URL and saves it to a file.
//
// Parameters:
//   - ctx: The context for cancellation
//   - url: The URL to download from
//   - filepath: The path to save the file
//
// Returns:
//   - error: Error if the download or save fails
func DownloadToFile(ctx context.Context, url, filepath string) error {
	reader, err := DownloadToReader(ctx, url)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err = io.Copy(out, reader); err != nil {
		return err
	}

	if err = out.Sync(); err != nil {
		return err
	}

	return nil
}

// CopyFile copies a file from source to destination.
// If the destination directory does not exist, it will be created.
//
// Parameters:
//   - src: The source file path
//   - dst: The destination file path
//
// Returns:
//   - error: Error if the copy fails
func CopyFile(src, dst string) error {
	if src == dst {
		return nil
	}
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	dstDir := filepath.Dir(dst)
	if _, err := os.Stat(dstDir); err == nil {
		// 存在
	} else if os.IsNotExist(err) {
		// 不存在，创建
		if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
			return err
		}
	} else {
		// 其他错误
		return err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err = dstFile.Sync(); err != nil {
		return err
	}
	if err := dstFile.Chmod(srcInfo.Mode()); err != nil {
		return err
	}
	return nil
}

// CopyDir copies an entire directory from source to destination.
//
// Parameters:
//   - src: The source directory path
//   - dst: The destination directory path
//
// Returns:
//   - error: Error if the copy fails
func CopyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, relPath), d.Type())
		} else {
			return CopyFile(path, filepath.Join(dst, relPath))
		}
	})
}

// Unzip extracts a zip archive to the specified destination directory.
//
// Parameters:
//   - src: The path to the zip file
//   - dst: The destination directory
//
// Returns:
//   - error: Error if extraction fails
func Unzip(src, dst string) (err error) {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, zipFile := range zipReader.File {
		if err := extractAndWriteFile(zipFile, dst); err != nil {
			return err
		}
	}
	return nil
}

func extractAndWriteFile(zipFile *zip.File, dst string) error {
	dstFilePath := filepath.Join(dst, zipFile.Name)

	if zipFile.FileInfo().IsDir() {
		return os.MkdirAll(dstFilePath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(dstFilePath), os.ModePerm); err != nil {
		return err
	}

	inFile, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return err
	}

	return outFile.Chmod(zipFile.FileInfo().Mode())
}

// IsDir reports whether the named file is a directory.
//
// Parameters:
//   - filepath: The path to check
//
// Returns:
//   - bool: True if the path is a directory, false otherwise
func IsDir(filepath string) bool {
	f, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// GetSize returns the size of the file at the given path.
//
// Parameters:
//   - path: The file path
//
// Returns:
//   - int64: The file size in bytes
//   - error: Error if the file does not exist
func GetSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

const (
	Byte     int64 = 1
	Kilobyte       = 1024 * Byte
	Megabyte       = 1024 * Kilobyte
	Gigabyte       = 1024 * Megabyte
	Terabyte       = 1024 * Gigabyte
	Petabyte       = 1024 * Terabyte
	Exabyte        = 1024 * Petabyte
	// Zettabyte          = 1024 * Exabyte
	// Yottabyte          = 1024 * Zettabyte
	// Brontobyte         = 1024 * Yottabyte
)

// HumanReadableSize converts a size in bytes to a human-readable string.
//
// Parameters:
//   - size: The size in bytes
//
// Returns:
//   - string: The human-readable size (e.g., "1.5GB")
func HumanReadableSize(size int64) string {
	s := size
	if s < 0 {
		s = -s
	}
	builder := stringx.Builder{}
	if s >= Exabyte {
		eb := s / Exabyte
		s = s % Exabyte
		_ = builder.WriteInt(eb, 10)
		_, _ = builder.WriteString("EB")
	}

	if s >= Petabyte {
		pb := s / Petabyte
		s = s % Petabyte
		_ = builder.WriteInt(pb, 10)
		_, _ = builder.WriteString("PB")
	}

	if s >= Terabyte {
		tb := s / Terabyte
		s = s % Terabyte
		_ = builder.WriteInt(tb, 10)
		_, _ = builder.WriteString("TB")
	}

	if s >= Gigabyte {
		gb := s / Gigabyte
		s = s % Gigabyte
		_ = builder.WriteInt(gb, 10)
		_, _ = builder.WriteString("GB")
	}

	if s >= Megabyte {
		mb := s / Megabyte
		s = s % Megabyte
		_ = builder.WriteInt(mb, 10)
		_, _ = builder.WriteString("MB")
	}

	if s >= Kilobyte {
		kb := s / Kilobyte
		s = s % Kilobyte
		_ = builder.WriteInt(kb, 10)
		_, _ = builder.WriteString("KB")
	}

	if s >= Byte {
		b := s / Byte
		// s = s % Byte
		_ = builder.WriteInt(b, 10)
		_, _ = builder.WriteString("B")
	}
	return builder.String()
}
