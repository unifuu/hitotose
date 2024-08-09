```
# Enable the fileinfo Extension
extension=fileinfo

# Create a New Laravel Project
composer create-project --prefer-dist laravel/laravel [PROJECT_NAME]

# Install Laravel MongoDB Package 
composer require mongodb/laravel-mongodb:^4.0
```

```
1. Download the MongoDB Extension:
https://pecl.php.net/package/mongodb

2. Place the downloaded file to `\php\ext`

3. Edit the `php.ini` and add the following line
extension=mongodb
```