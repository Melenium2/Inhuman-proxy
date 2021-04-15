package storage

import "context"

func (r RedisStorage) ChangeBlockStatus(ctx context.Context, address string, status string) error {
	return r.changeBlockStatus(ctx, address, status)
}

func (r RedisStorage) Delete(ctx context.Context, address string) error {
	return r.delete(ctx, address)
}
