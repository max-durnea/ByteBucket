1. Implement user authentification and storing in Postgresql database
    - Create the table (id, username, email, password_hash, created_at, updated_at)
    - Fix migration bash scripts
2. Implement authorization using middleware for access token check
3. Add database refresh token cleanup
4. Add database cleanup endpoint for dev
    - use an environmental variable to check the platform
5. Add file uploading and downloading ability
6. Add other features
