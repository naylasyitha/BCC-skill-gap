CREATE TABLE users (
    id UUID PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(500) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    education_level VARCHAR(50),
    major VARCHAR(255),
    institution VARCHAR(255),
    graduation_year INT,
    is_premium BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    refresh_token VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE skills (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE careers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE career_skills (
    id UUID PRIMARY KEY,
    career_id UUID NOT NULL,
    skill_id UUID NOT NULL,
    priority INT NOT NULL,
    required_level VARCHAR(50) NOT NULL,
    CONSTRAINT fk_career_skills_career FOREIGN KEY (career_id) REFERENCES careers(id) ON DELETE CASCADE,
    CONSTRAINT fk_career_skills_skill FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE questions (
    id UUID PRIMARY KEY,
    skill_id UUID NOT NULL,
    level VARCHAR(50) NOT NULL,
    question_content TEXT NOT NULL,
    option_a VARCHAR(255) NOT NULL,
    option_b VARCHAR(255) NOT NULL,
    option_c VARCHAR(255) NOT NULL,
    option_d VARCHAR(255) NOT NULL,
    answer CHAR(1) NOT NULL,
    explanation TEXT,
    CONSTRAINT fk_questions_skill FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE materials (
    id UUID PRIMARY KEY,
    skill_id UUID NOT NULL,
    level VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    video_url VARCHAR(500),
    order_number INT NOT NULL,
    CONSTRAINT fk_materials_skill FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE user_career_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    career_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'on_process',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_user_career_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_career_sessions_career FOREIGN KEY (career_id) REFERENCES careers(id) ON DELETE CASCADE
);

CREATE TABLE self_assessment_skills (
    id UUID PRIMARY KEY,
    user_career_session_id UUID NOT NULL,
    skill_id UUID NOT NULL,
    user_level VARCHAR(50),
    user_final_level VARCHAR(50),
    quiz_score INT DEFAULT 0,
    CONSTRAINT fk_sas_session FOREIGN KEY (user_career_session_id) REFERENCES user_career_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_sas_skill FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE quiz_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    user_career_session_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'on_process',
    score FLOAT8,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_quiz_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_quiz_sessions_session FOREIGN KEY (user_career_session_id) REFERENCES user_career_sessions(id) ON DELETE CASCADE
);

CREATE TABLE learning_path_progresses (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    user_career_session_id UUID NOT NULL,
    material_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'not_started',
    completed_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_lpp_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_lpp_session FOREIGN KEY (user_career_session_id) REFERENCES user_career_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_lpp_material FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
);

CREATE TABLE quiz_answers (
    id UUID PRIMARY KEY,
    quiz_session_id UUID NOT NULL,
    question_id UUID NOT NULL,
    user_answer CHAR(1),
    is_correct BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_quiz_answers_session FOREIGN KEY (quiz_session_id) REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_quiz_answers_question FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

CREATE TABLE transactions (
    id VARCHAR(50) PRIMARY KEY,
    user_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    snap_token VARCHAR(255),
    snap_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);