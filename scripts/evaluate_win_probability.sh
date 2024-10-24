#!/bin/bash
( 
echo "setoption name Skill Level value $2" ;
echo "position fen $3" ;
echo "go depth 20" ;  # Увеличиваем глубину для более точной оценки
sleep 1
) | $1 | grep -A 10 "info" | grep "score"  # Фильтруем вывод для получения оценки
