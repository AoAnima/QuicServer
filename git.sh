#!/bin/bash
cd /home/exmaao@mfsk.int/Файлы/AoAnima.ru/QuicMarket
git add -A
if [ -n "$(git status --porcelain)" ]; then
    git commit -m "Автоматическая отправка"
    git push
fi

cd /home/exmaao@mfsk.int/Файлы/AoAnima.ru/jet
git add -A
if [ -n "$(git status --porcelain)" ]; then
    git commit -m "Автоматическая отправка"
    git push
fi


cd /home/exmaao@mfsk.int/Файлы/AoAnima.ru/QuicMarket/Logger
git add -A
if [ -n "$(git status --porcelain)" ]; then
    git commit -m "Автоматическая отправка"
    git push
fi