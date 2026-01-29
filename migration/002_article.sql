-- Create articles table
CREATE TABLE IF NOT EXISTS articles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_articles_author_id (author_id),
    INDEX idx_articles_created_at (created_at),
    INDEX idx_articles_title (title),
    
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);