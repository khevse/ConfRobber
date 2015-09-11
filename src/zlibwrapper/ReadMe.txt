Источники:

1. http://www.swig.org/Doc3.0/SWIGDocumentation.html
2. http://www.swig.org/Doc3.0/Go.html
3. http://zacg.github.io/blog/2013/06/06/calling-c-plus-plus-code-from-go-with-swig/
4. http://akrennmair.github.io/golang-cgo-slides/#9
5. https://github.com/golang/go/wiki/cgo
6. https://talks.golang.org/2015/state-of-go-may.slide#23
7. http://blog.xebia.com/2014/07/04/create-the-smallest-possible-docker-container/
8. http://www.ibm.com/developerworks/ru/library/au-swig/

Перед началом сборки необходимо скачать:
    
- исходники zlib: http://www.zlib.net/
- swig:http://www.swig.org/survey.html
- mingw64:http://sourceforge.net/projects/mingw-w64/
- cmake: http://www.cmake.org/download/

После установки выше описанных программы необходимо внести изменения в файл
"compile.bat" и запустить его на испольнение с одним из параметров:
- clean   - очистить текущую директорию от временных файлов
- release - сборать с удалением временных файлов (режим по умолчанию)