## start backend server
```
go run .
```

## build docker image
```
docker build -t minpeter/hijackjwtadmin .
```

## run docker images
```
docker run --rm -p 4000:4000 minpeter/hijackjwtadmin
```

## 유출된 wordlist 생성 방법
```
crunch 8 8 qwertyuiopasdfghjklzxcvbnm0987654321 -t cr34m@@@ -o ./wordlist.txt
```


## 주의사항
문제 파일 배포할때 꼭 production.env, solution.md 파일을 제거하고 압축하여 배포하세요!!  

## 참고사항
문제에서 유출된 패스워드 리스트는 용량상의 문제로 업로드하지 않았습니다.  
위의 명령어를 활용해 유출된 패스워드를 생성해보세요!!  
`rest.http` 파일의 경우 vscode의 REST Client 익스텐션을 이용해 사용할 수 있습니다.  
