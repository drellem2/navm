nasm.lnk -f win64 C:\dev\navm\out.asm -o C:\dev\navm\test.o
ld test.o -o assembly.exe
assembly.exe
echo %ERRORLEVEL%
