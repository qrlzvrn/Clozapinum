#!/bin/bash

if [ -f .env/db.env ]
then
	echo "Файл конфигурации для базы данных уже существует"
else
	echo "Сконфигурируем базу данных"

    dbusername_str="POSTGRES_USER="
	read -p "Введите имя пользователя базы данных: " dbusername




	ask_password() {
		dbpassword_str="POSTGRES_PASSWORD="
		read -s -p "Введите пароль пользователя базы данных: " dbpassword
        echo ""

		read -s -p "Повторите пароль пользователя базы данных: " dbpassword_check
        echo ""

		if [[ "$dbpassword" != "$dbpassword_check" ]]
		then
			if [[ "10#$dbpassword" -ne "10#$dbpassword_check" ]]
			then
				echo "Ваши пароли не совпадают"
				ask_password
			fi
			echo "Ваши пароли не совпадают"
			ask_password
		fi
	}

	ask_password

	dbname_str="POSTGRES_DB="
	read -p "Введите название базы данных: " dbname
	

	if [ -d .env ]
	then
		touch .env/db.env

		echo "$dbpassword_str$dbpassword" > .env/db.env
		echo "$dbusername_str$dbusername" >> .env/db.env
		echo "$dbname_str$dbname" >> .env/db.env
	else
		mkdir .env

		touch .env/db.env

		echo "$dbpassword_str$dbpassword" > .env/db.env
		echo "$dbusername_str$dbusername" >> .env/db.env
		echo "$dbname_str$dbname" >> .env/db.env
	fi
fi


if [ -f .env/tgBot.env ]
then
	echo "файл конфигурации бота уже существует"
else
	token_str="TELEGRAM_APITOKEN="
	read -p "Введите токен бота: " token

	fullchain_str="SSL_FULLCHAIN="
	read -e -p "Укажите путь к fullchain.pem: " fullchain

	privkey_str="SSL_PRIVKEY="
	read -e -p "Укажите путь к privkey.pem: " privkey




	db_host_str="DB_HOST="
	read -p "Введите адрес или домен базы данных: " host

	db_port_str="DB_PORT="
	read -p "Введите порт базы данных: " port

	db_username_str="DB_USERNAME="
	db_password_str="DB_PASSWORD="
	db_name_str="DB_NAME="


	touch .env/tgBot.env

	echo "$token_str$token" > .env/tgBot.env
	echo "$fullchain_str$fullchain" >> .env/tgBot.env
	echo "$privkey_str$privkey" >> .env/tgBot.env
	echo "$db_host_str$host" >> .env/tgBot.env
	echo "$db_port_str$port" >> .env/tgBot.env


	echo "$db_username_str$dbusername" >> .env/tgBot.env
	echo "$db_password_str$dbpassword" >> .env/tgBot.env
	echo "$db_name_str$dbname" >> .env/tgBot.env

fi