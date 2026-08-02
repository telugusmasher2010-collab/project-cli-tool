package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	apperrors "github.com/telugusmasher2010-collab/project-cli-tool/internal/errors"
)

// untar extracts a gzipped tar archive into dir.
func untar(src, dir string) error {
	f, err := os.Open(src)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to open archive", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to read gzip stream", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return apperrors.Wrap(apperrors.ErrFilesystem, "failed to read tar stream", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		name := filepath.Base(hdr.Name)
		if strings.HasPrefix(name, ".") {
			continue
		}
		target := filepath.Join(dir, name)
		if err := copyFile(target, tr, hdr.FileInfo().Mode()); err != nil {
			return err
		}
	}
	return nil
}

// unzip extracts a zip archive into dir.
func unzip(src, dir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to open archive", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := filepath.Base(f.Name)
		if strings.HasPrefix(name, ".") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return apperrors.Wrap(apperrors.ErrFilesystem, "failed to open archived entry", err)
		}
		err = copyFile(filepath.Join(dir, name), rc, f.Mode())
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// copyFile writes src into target with the given mode.
func copyFile(target string, src io.Reader, mode os.FileMode) error {
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to create file", err)
	}
	if _, err := io.Copy(out, src); err != nil {
		out.Close()
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to write file", err)
	}
	if err := out.Close(); err != nil {
		return apperrors.Wrap(apperrors.ErrFilesystem, "failed to close file", err)
	}
	return nil
}
