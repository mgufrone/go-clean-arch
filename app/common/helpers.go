package common

type TryError func() error

func Try(calls ...TryError) error {
	for _, c := range calls {
		if err := c(); err != nil {
			return err
		}
	}
	return nil
}

