[LND](https://github.com/lightningnetwork/lnd),[whatsat](https://github.com/joostjager/whatsat)을 활용하여 통신을 구현하였다.

Whatsat이란 라이트닝 네트워크를 활용한 채팅 프로그램 이다.

라이트닝 네트워크에서 송금할때 메세지를 추가해 마치 채팅을 하는것처럼 보이게 해준다.

본 논문에서는 Whatsat을 이용해 메세지 부분에 명령어를 넣어 사물인터넷 기기가 받게되면 해당 명령을 실행하게 하였다.

-Whatsat을 이용한 기본 채팅 

![라이트닝-네트워크-채팅](https://user-images.githubusercontent.com/33947681/88308039-8bcc1000-cd47-11ea-871d-6b9b078801fa.gif)


-whatsat의 메세지를 수정하여 라즈베리파이에게 온도를 요청하고 받는 영상
![무제](https://user-images.githubusercontent.com/33947681/88308175-b027ec80-cd47-11ea-9246-a8e0bc7adc91.gif)

-노트북에서 온도를 요청하면 라즈베리파이에 부착된 센서가 온도를 측정하고 다시 노트북으로 전송하는 영상
![ezgif com-optimize](https://user-images.githubusercontent.com/33947681/88310412-9b992380-cd4a-11ea-87db-e5893dd31a24.gif)

