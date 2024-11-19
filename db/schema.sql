-- Should be able to adapt this to any one table csv file to work with the program
CREATE TABLE dogs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    origin VARCHAR(100),
    type VARCHAR(50),
    unique_feature TEXT,
    friendly_rating SMALLINT CHECK (friendly_rating BETWEEN 1 AND 10),
    life_span SMALLINT CHECK (life_span > 0),
    size VARCHAR(50),
    grooming_needs VARCHAR(100),
    exercise_requirements NUMERIC(3,1) CHECK (exercise_requirements >= 0),
    good_with_children VARCHAR(50),
    intelligence_rating SMALLINT CHECK (intelligence_rating BETWEEN 1 AND 10),
    shedding_level VARCHAR(50),
    health_issues_risk VARCHAR(50),
    average_weight NUMERIC(5,2),
    training_difficulty SMALLINT CHECK (training_difficulty BETWEEN 1 AND 10)
);

-- Due to read only nature of project we are gonna learn about indexing 
-- Should have an initial performance cost i believe but should be worth it 
CREATE INDEX idx_dogs_name ON dogs(name);
CREATE INDEX idx_dogs_origin ON dogs(origin);
CREATE INDEX idx_dogs_type ON dogs(type);
CREATE INDEX idx_dogs_unique_feature ON dogs(unique_feature);
CREATE INDEX idx_dogs_friendly_rating ON dogs(friendly_rating);
CREATE INDEX idx_dogs_life_span ON dogs(life_span);
CREATE INDEX idx_dogs_size ON dogs(size);
CREATE INDEX idx_dogs_grooming_needs ON dogs(grooming_needs);
CREATE INDEX idx_dogs_exercise_requirements ON dogs(exercise_requirements);
CREATE INDEX idx_dogs_good_with_children ON dogs(good_with_children);
CREATE INDEX idx_dogs_intelligence_rating ON dogs(intelligence_rating);
CREATE INDEX idx_dogs_shedding_level ON dogs(shedding_level);
CREATE INDEX idx_dogs_health_issues_risk ON dogs(health_issues_risk);
CREATE INDEX idx_dogs_average_weight ON dogs(average_weight);
CREATE INDEX idx_dogs_training_difficulty ON dogs(training_difficulty);

-- Make Sure the headers in the csv file match exactly to your schema file 
COPY dogs(
    name,
    origin,
    type,
    unique_feature,
    friendly_rating,
    life_span,
    size,
    grooming_needs,
    exercise_requirements,
    good_with_children,
    intelligence_rating,
    shedding_level,
    health_issues_risk,
    average_weight,
    training_difficulty
) FROM '/data/DogData.csv' DELIMITER ',' CSV HEADER;
