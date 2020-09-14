%⭐️⭐️%

블록체인으로 사물 인터넷의 무결성을 보장하는 시스템

A system to ensure the integrity of Internet of things by using Blockchain


김영석

Young-Seok Kim

youngseokaaa@gmail.com


사물 인터넷 기기의 개수가 늘어나는 만큼 보안위협 또한 같이 늘어나고 있다. IBM의 ‘The Dangers of Smart City Hacking’[5]에 따르면 사물 인터넷 기기의 데이터를 조작하는 행위가 일어난다면 심각한 문제를 일으킬 수 있다고 경고한다. 이에 본 논문에서는 사물 인터넷 시스템의 무결성을 보장하기 위해 비트코인과 라이트닝 네트워크를 사용한다. 블록체인 기술은 각각의 사용자가 데이터를 저장할 때 특정 합의 알고리즘을 통해 블록을 생성하는 방식으로 데이터를 사용자들과 공유하고 저장한다. 이는 탈 중앙적인 구조로 누구나 참여하여 데이터를 저장할 수 있지만 악의적인 사용자 또한 참여하여 시스템을 무력화시킬 수 있다. 블록체인에서는 이를 흔히 ‘51% 공격’ 이라고 부른다. 본 논문에서는 정상적인 사용자들에게 51% 공격을 했을 때 무결성이 지켜지는지 확인 하기 위한 실험을 진행하였다. 그 결과 악의적인 사용자가 시스템의 절반 이상의 해시 연산력을 보유할 경우 데이터를 수정해 무결성을 위협할 수 있다는 점을 확인했다. 하지만 국내에서 진행하는 ‘스마트홈 10만 가구 구축’ 프로젝트에 적용했을때 10만 가구 이상의 해시 연산력을 보유하여야 무결성을 위협할 수 있기에 현실적으로 불가능 하다. 그래서 본 논문에서는 사물 인터넷의 무결성을 보장하기 위해 비트코인을 사용한 시스템을 제안한다.





실험을 위해 기존의 비트코인에 하드포크를 진행해 영석코인을 만들어 진행했다.

lightningnetwork_connect는 사용자와 사물 인터넷 간의 통신을 위해 만든 어플리케이션이다.

LND를 기반으로한 whatsat을 수정하여 사용하였다.
