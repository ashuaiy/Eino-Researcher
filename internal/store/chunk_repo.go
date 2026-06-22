package store

import (
	"context"
	"fmt"

	"github.com/pgvector/pgvector-go"

	"eino-researcher/internal/model"
)

type ChunkRepository interface {
	CreateMany(ctx context.Context, chunks []model.Chunk) error
	Search(ctx context.Context, embedding []float32, topK int) ([]model.Evidence, error)
}

type NoopChunkRepository struct{}

func (r NoopChunkRepository) CreateMany(ctx context.Context, chunks []model.Chunk) error {
	// TODO: insert chunk rows with pgvector embeddings.
	return nil
}

func (r NoopChunkRepository) Search(ctx context.Context, embedding []float32, topK int) ([]model.Evidence, error) {
	// TODO: implement vector similarity search with pgvector.
	return []model.Evidence{}, nil
}

type PostgresChunkRepository struct {
	db  DBTX
	dim int
}

func NewPostgresChunkRepository(db DBTX, dim int) *PostgresChunkRepository {
	return &PostgresChunkRepository{db: db, dim: dim}
}

func (r *PostgresChunkRepository) CreateMany(ctx context.Context, chunks []model.Chunk) error {
	for _, chunk := range chunks {
		if r.dim > 0 && len(chunk.Embedding) != r.dim {
			return fmt.Errorf(
				"chunk embedding dimension mismatch: expected %d, got %d",
				r.dim,
				len(chunk.Embedding),
			)
		}
		_, err := r.db.Exec(
			ctx,
			`INSERT INTO chunks
			 (id, document_id, content, chunk_index, token_count, embedding, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			chunk.ID,
			chunk.DocumentID,
			chunk.Content,
			chunk.Index,
			chunk.TokenCount,
			pgvector.NewVector(chunk.Embedding),
			chunk.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert chunk %d: %w", chunk.Index, err)
		}
	}
	return nil
}

func (r *PostgresChunkRepository) Search(ctx context.Context, embedding []float32, topK int) ([]model.Evidence, error) {
	if r.dim > 0 && len(embedding) != r.dim {
		return nil, fmt.Errorf(
			"query embedding dimension mismatch: expected %d, got %d",
			r.dim,
			len(embedding),
		)
	}

	rows, err := r.db.Query(
		ctx,
		`SELECT
			c.id::text,
			c.document_id::text,
			d.title,
			COALESCE(d.source, ''),
			c.content,
			1 - (c.embedding <=> $1) AS score
		 FROM chunks c
		 JOIN documents d ON d.id = c.document_id
		 WHERE c.embedding IS NOT NULL
		 ORDER BY c.embedding <=> $1
		 LIMIT $2`,
		pgvector.NewVector(embedding),
		topK,
	)
	if err != nil {
		return nil, fmt.Errorf("search chunks: %w", err)
	}
	defer rows.Close()

	results := make([]model.Evidence, 0, topK)
	for rows.Next() {
		var evidence model.Evidence
		if err := rows.Scan(
			&evidence.ID,
			&evidence.DocumentID,
			&evidence.Title,
			&evidence.Source,
			&evidence.Content,
			&evidence.Score,
		); err != nil {
			return nil, fmt.Errorf("scan chunk search result: %w", err)
		}
		evidence.SourceType = "local_doc"
		results = append(results, evidence)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chunk search results: %w", err)
	}
	return results, nil
}
