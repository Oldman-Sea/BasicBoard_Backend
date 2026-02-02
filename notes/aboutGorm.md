## gORM

### ORM (Object Relation Mapping)

- 객체 지향 프로그래밍 언어(Java, Python 등)의 객체와 관계형 데이터베이스(RDB)의 테이블을 매핑하여, SQL 쿼리 없이 객체 중심으로 데이터를 다루는 기술.

- 장점
  1. 객체 중심 개발: SQL 문 대신 프로그래밍 언어의 메서드(예: save(), find())로 DB 조작.
  2. 생산성 향상: SQL 쿼리 반복 작성 감소 및 가독성 향상.
  3. DBMS 종속성 감소: SQL을 ORM이 자동 생성하므로 다른 DB로 전환 용이.
  4. 재사용 및 유지보수: 객체 모델링 중심의 설계로 코드 재사용성 증가.

### gORM

- gORM은 golang에서 사용 가능한 ORM 라이브러리.
- gORM 설치: go get -u gorm.io/gorm
- gORM용 MySQL 드라이버 설치: go get -u github.com/go-sql-driver/mysql

### First vs Find vs Scan

- 공통점: 셋 다 DB 결과를 Go 변수에 집어넣는 역할을 함.
- 차이점: **어떤 형태의 그릇에, 어떻게 채우느냐**

- First
  조건에 맞는 **첫번째 레코드**를 **struct**에 채움.

  ```go
  var post models.Post
  err := db.Where("id = ?", 10).First(&post).Error
  // sql문 번역
  // SELECT * FROM posts WHERE id = 10 LIMIT 1;
  ```

  - struct(해당 코드에서는 models.Post) 전체 컬럼 채움.
  - 결과 없으면 → record not found 에러
  - 단일 row 용
  - **엔티티 하나 가져올 때 사용**

- Find
  조건에 맞는 **여러 레코드를**를 **slice**에 채움.

  ```go
  var posts []models.Post
  err := db.Where("user_id = ?", 1).Find(&posts).Error
  // sql문 번역
  // SELECT * FROM posts WHERE user_id = 1;
  ```

  - 여러 row 가능
  - 결과 없어도 에러 안남.(빈 slice 반환)
  - 목록 조회용
  - **리스트 페이지**

- Scan
  쿼리 결과를 **임의의 변수/struct**에 그대로 복사

  ```go
  var id uint
  db.Select("id").Scan(&id)
  // sql문 번역
  // SELECT id FROM posts ...
  ```

  - Model struct가 필요 없음.
  - 컬럼 수 = 변수/필드 수만 맞으면 됨.
  - 부분 조회/집계/커스텀 쿼리에 최적
  - **SQL 결과를 그냥 복사해 반환**

| 구분         | First  | Find     | Scan             |
| ------------ | ------ | -------- | ---------------- |
| 대상         | struct | slice    | 아무 변수        |
| 컬럼         | 전체   | 전체     | 선택한 것만      |
| row 수       | 1      | N        | 자유             |
| 결과 없을 때 | 에러   | 빈 slice | 값 그대로        |
| 용도         | 상세   | 목록     | id만, count, sum |
