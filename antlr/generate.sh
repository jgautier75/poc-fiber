#!/bin/sh

antlr='java -jar antlr-4.13.2-complete.jar'

$antlr -Dlanguage=Go Filter.g4 -o parser
$antlr -Dlanguage=Go Filter.g4 -visitor -o parser
