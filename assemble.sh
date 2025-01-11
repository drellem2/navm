as out.as -o out.o && \
ld out.o -o a.out -l System -syslibroot $(xcrun -sdk macosx --show-sdk-path)  -e _start -arch arm64 && \
chmod +x a.out && \
./a.out
