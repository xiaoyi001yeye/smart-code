CREATE TABLE IF NOT EXISTS tasks (
    task_id TEXT PRIMARY KEY,
    input_path TEXT NOT NULL,
    result_path TEXT NOT NULL,
    code_language TEXT NOT NULL,
    task_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    current_step INT,
    is_completed BOOLEAN
);

-- CREATE TABLE IF NOT EXISTS repositories (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     branch_name VARCHAR(255),
--     repo_path TEXT NOT NULL,
--     username VARCHAR(255) NOT NULL,
--     password TEXT NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );


-- INSERT INTO repositories (name, branch_name, repo_path, username, password) VALUES('Example Repository', 'main', 'https://github.com/spring-projects/spring-petclinic', 'user', 'password') ON CONFLICT (id) DO NOTHING;