package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Dog struct {
	ID                   int
	Name                 string
	Origin               string
	Type                 string
	UniqueFeature        string
	FriendlyRating       int16
	LifeSpan             int16
	Size                 string
	GroomingNeeds        string
	ExerciseRequirements float64
	GoodWithChildren     string
	IntelligenceRating   int16
	SheddingLevel        string
	HealthIssuesRisk     string
	AverageWeight        float64
	TrainingDifficulty   int16
}

func NewPostgreSQLClient(dbURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, err
}

func DogBreedsPaginated(dbClient *sql.DB, page int, limit int) ([]Dog, error){
    if limit > 100 {
        limit = 100
    }

    offset := (page - 1) * limit

    query := `
        SELECT * FROM dogs LIMIT $1 OFFSET $2
    `

    rows, err := dbClient.Query(query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("error executing query: %v", err)
    }

    var dogs []Dog

    for rows.Next() {
		var dog Dog
		if err := rows.Scan(
			&dog.ID, &dog.Name, &dog.Origin, &dog.Type, &dog.UniqueFeature, 
			&dog.FriendlyRating, &dog.LifeSpan, &dog.Size, &dog.GroomingNeeds, 
			&dog.ExerciseRequirements, &dog.GoodWithChildren, &dog.IntelligenceRating, 
			&dog.SheddingLevel, &dog.HealthIssuesRisk, &dog.AverageWeight, &dog.TrainingDifficulty,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		dogs = append(dogs, dog)
	}

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error with ros: %v", err)
    }

    return dogs, nil

}