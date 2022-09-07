# 누출된 패스워드 목록과 jwt 관리자 세션 탈취

## 시나리오

당신은 보안회사 hsocSec의 화이트해커이다.

의류 쇼핑몰 회사인 **Cream**에서 보안관리자의 실수로 인터넷에 회사 시스템에서 이용되는 암호 리스트가 불특정 다수에게 유출되는 사건이 발생하였다.

유출 사실을 파악한 직후 보안관리자가 빠르게 파일을 삭제하긴하였지만 이로인해 발생할 수 있는 보안취약점을 파악하고 해결하기 위해 HsocSec에 보안 점검을 의뢰하였다.

시스템 전반에 거쳐 사용되는 암호리스트가 유출된만큼 전부 점검이 필요하지만 어느 시스템에 암호를 변경해야될지 모르는 상황에서 가장 중요한 Cream 회사의 인증서버를 점검하기로 결정했다.

다음은 Cream의 인증 서버의 주소이다.

로그인 로직에 유출된 패스워드로 인해 발생할 수 있는 피해를 조사하며 관리자 계정에 존재하는 플래그를 찾아보자.