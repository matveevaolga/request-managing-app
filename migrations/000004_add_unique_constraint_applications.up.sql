DELETE FROM applications a1 USING applications a2 
WHERE a1.id < a2.id 
  AND a1.project_name = a2.project_name 
  AND a1.email = a2.email;

ALTER TABLE applications ADD CONSTRAINT unique_project_email UNIQUE (project_name, email);