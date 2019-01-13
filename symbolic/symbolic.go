package symbolic

import (
	"fmt"
	"os"
)

func Ln(source string, target string) error {
	err := Unlink(target)
	if err != nil {
		return err
	}

	err = os.Symlink(source, target)
	if err != nil {
		return err
	}

	return nil
}

func Unlink(target string) error {
	if _, err := os.Lstat(target); err == nil {
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("failed to unlink: %+v", err)
		}

		return nil
	} else if os.IsNotExist(err) {
		// no link exists at target
		return nil
	}

	return fmt.Errorf("unknown error checking symlink. SHOULD NEVER GET HERE...")
}

func Readlink(target string) (string, error) {
	if _, err := os.Lstat(target); err == nil {
		if linkPath, err := os.Readlink(target); err == nil {
			return linkPath, nil
		} else {
			return "", fmt.Errorf("Error: unable to resolve symlink %s\n\n%s\n", target, err)
		}
	} else if os.IsNotExist(err) {
		return "", err
	}

	return "", fmt.Errorf("Error: unknown error. SHOULD NEVER BE ABLE TO GET HERE")
}