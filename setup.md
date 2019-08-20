env GOOS=linux go build -o=jiva-picsum ./main.go


```bash
export CGO_CFLAGS_ALLOW=-Xpreprocessor; 
export PKG_CONFIG_PATH=/usr/local/opt/libffi/lib/pkgconfig; 
go get -u github.com/davidbyttow/govips/pkg/vips
```


MySQL database setup

```mysql
CREATE DATABASE `imaginary`;
USE imaginary;
CREATE TABLE `image` (
  `id` varchar(255) NOT NULL,
  `author` varchar(30) NOT NULL,
  `width` int(11) NOT NULL,
  `height` int(11) NOT NULL,
  `url` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

## Configuration file setup

conf.yaml

```yaml

```