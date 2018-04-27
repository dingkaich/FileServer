create table fileinfo (
    filename varchar(1024) not null CONSTRAINT filename_pk PRIMARY KEY,
    filemd5  varchar(512),
    fileserverpath varchar(1024),
    serverIp    varchar(25),
    clientIp varchar(25)
    uploadtime date
);


create table userinfo(
    username varchar(64),
	passwd   varchar(128)
);