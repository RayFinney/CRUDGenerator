CREATE TABLE news_articles
(
	uuid UUID NOT NULL,
	title VARCHAR(50) NOT NULL,
	article TEXT NOT NULL,
	views INT NULL,
	rating REAL NOT NULL,
	public BOOL NOT NULL,
	created_by UUID NOT NULL,
	modified_by UUID NULL,
	deleted_by UUID NULL,
	created_date TIMESTAMP NOT NULL,
	last_modified TIMESTAMP NULL,
	deleted_date TIMESTAMP NULL,
	PRIMARY KEY (uuid)
);