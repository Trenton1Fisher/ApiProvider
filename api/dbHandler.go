package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

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

//Whitelisted url params for filtered api endpoint
var DogTableColumns = []string{
    "id", "name", "origin", "type", "unique_feature", "friendly_rating",
    "life_span", "size", "grooming_needs", "exercise_requirements",
    "good_with_children", "intelligence_rating", "shedding_level",
    "health_issues_risk", "average_weight", "training_difficulty",
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

func DogById(dbClient *sql.DB, id int) (Dog, error) {
    query := `
        SELECT * FROM dogs WHERE id = $1
    `
    var dog Dog
    row := dbClient.QueryRow(query, id)
    err := row.Scan(
        &dog.ID, &dog.Name, &dog.Origin, &dog.Type, &dog.UniqueFeature, 
        &dog.FriendlyRating, &dog.LifeSpan, &dog.Size, &dog.GroomingNeeds, 
        &dog.ExerciseRequirements, &dog.GoodWithChildren, &dog.IntelligenceRating, 
        &dog.SheddingLevel, &dog.HealthIssuesRisk, &dog.AverageWeight, &dog.TrainingDifficulty,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return Dog{}, fmt.Errorf("no dog found with id %d", id)
        }
        return Dog{}, fmt.Errorf("error scanning row: %v", err)
    }

    return dog, nil
}

func DogsByFilter(dbClient *sql.DB, query string, values []interface{})([]Dog, error) {

   var dogs []Dog
   rows, err := dbClient.Query(query, values...)
   if err != nil {
       return nil, fmt.Errorf("error executing query: %v", err)
   }
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

func DogsByFilterQueryBuilder(params url.Values) (string, []interface{}) {
	values := []interface{}{}
	conditions := []string{}
	query := "SELECT * FROM dogs"

	for _, v := range DogTableColumns {
		paramValue := params.Get(v)
		if paramValue != "" {
			values = append(values, paramValue)
			conditions = append(conditions, fmt.Sprintf("%s = $%d", v, len(values)))
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, values
}
