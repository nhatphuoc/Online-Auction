# WNC - Final Project - Online Auction

Xây dựng ứng dụng Sàn Đấu Giá Trực Tuyến, gồm các phân hệ & chức năng sau

## 1. Phân hệ người dùng nặc danh - guest

### 1.1 Hệ thống Menu

* Hiển thị danh sách danh mục category
* Có 2 cấp danh mục

  * Điện tử ➠ Điện thoại di động
  * Điện tử ➠ Máy tính xách tay
  * Thời trang ➠ Giày
  * Thời trang ➠ Đồng hồ
  * …

### 1.2 Trang chủ

* Top 5 sản phẩm gần kết thúc
* Top 5 sản phẩm có nhiều lượt ra giá nhất
* Top 5 sản phẩm có giá cao nhất

### 1.3 Xem danh sách sản phẩm

* Theo danh mục category
* Có phân trang

### 1.4 Tìm kiếm sản phẩm

* Sử dụng kỹ thuật Full-text search, cho phép tìm kiếm tiếng Việt không dấu
* Tìm theo tên sản phẩm and/or tìm theo danh mục
* Phân trang kết quả
* Sắp xếp theo ý người dùng

  * Thời gian kết thúc giảm dần
  * Giá tăng dần
* Những sản phẩm mới đăng (trong vòng N phút) sẽ thể hiện nổi bật

#### 1.4.1 Sản phẩm hiển thị trên trang danh sách gồm:

* Ảnh đại diện sản phẩm
* Tên sản phẩm
* Giá hiện tại
* Thông tin bidder đang đặt giá cao nhất
* Giá mua ngay (nếu có)
* Ngày đăng sản phẩm
* Thời gian còn lại
* Số lượt ra giá hiện tại
* Người dùng có thể click vào category để chuyển nhanh sang màn hình xem danh sách sản phẩm

### 1.5 Xem chi tiết sản phẩm

* Nội dung đầy đủ của sản phẩm
* Ảnh đại diện (size lớn)
* Các ảnh phụ (ít nhất 3 ảnh)
* Tên sản phẩm
* Giá hiện tại
* Giá mua ngay (nếu có)
* Thông tin người bán & điểm đánh giá
* Thông tin người đặt giá cao nhất hiện tại & điểm đánh giá
* Thời điểm đăng
* Thời điểm kết thúc
* Nếu thời điểm kết thúc < 3 ngày thì hiển thị relative time
* Mô tả chi tiết sản phẩm
* Lịch sử các câu hỏi và câu trả lời
* 5 sản phẩm khác cùng chuyên mục

### 1.6 Đăng ký

* Người dùng cần đăng ký để có thể đặt giá
* reCaptcha
* Mật khẩu mã hoá bằng bcrypt hoặc scrypt
* Thông tin yêu cầu:

  * Họ tên
  * Địa chỉ
  * Email (không trùng)
  * Xác nhận OTP

---

## 2. Phân hệ người mua bidder

### 2.1 Watch List

* Lưu 1 sản phẩm vào danh sách yêu thích tại danh sách sản phẩm hoặc chi tiết sản phẩm

### 2.2 Ra giá

* Chỉ bidder có điểm đánh giá >= 80% mới được ra giá
* Bidder chưa từng được đánh giá có thể ra giá nếu người bán cho phép
* Hệ thống đề nghị giá hợp lệ (giá hiện tại + bước giá)
* Yêu cầu xác nhận trước khi gửi giá

### 2.3 Xem lịch sử đấu giá

* Thông tin người ra giá được che 1 phần

### 2.4 Hỏi người bán

* Thực hiện tại trang chi tiết sản phẩm
* Người bán nhận email thông báo và link trả lời

### 2.5 Quản lý hồ sơ cá nhân

* Đổi email, họ tên, mật khẩu
* Xem điểm đánh giá và chi tiết các đánh giá
* Xem danh sách yêu thích
* Xem danh sách sản phẩm đang đấu giá
* Xem danh sách sản phẩm đã thắng
* Được phép đánh giá người bán (+1 / -1) kèm nhận xét

### 2.6 Xin được bán trong 7 ngày

* Bidder gửi yêu cầu nâng cấp thành seller
* Admin duyệt

---

## 3. Người bán - seller

### 3.1 Đăng sản phẩm đấu giá

* Tên sản phẩm
* Ít nhất 3 ảnh
* Giá khởi điểm
* Bước giá
* Giá mua ngay
* Mô tả sản phẩm (WYSIWYG)
* Tự động gia hạn hay không
* Có tham số 5 phút và 10 phút (admin chỉnh)

### 3.2 Bổ sung thông tin mô tả

* Chỉ được append, không thay thế nội dung cũ

```
電源入り撮影出来ましたが細部の機能までは確認していません。
不得意ジャンルの買い取り品の為細かい確認出来る知識がありません、ご了承ください。
簡単な確認方法が有れば確認しますので方法等質問欄からお願いします、終了日の質問には答えられない場合があります。
付属品、状態は画像でご確認ください。
当方詳しくありませんので高度な質問には答えられない場合がありますがご了承ください。
発送は佐川急便元払いを予定しています、破損防止の為梱包サイズが大きくなる事がありますがご了承下さい。
中古品の為NC/NRでお願いします。

✏️ 31/10/2025

- が大きくなる事がありますがご了承下さい。

✏️ 5/11/2025

- 不得意ジャンルの買い取り品の為細かい確認出来る知識がありません、ご了承ください。
```

### 3.3 Từ chối lượt ra giá

* Người bị từ chối không được phép đấu giá sản phẩm đó
* Nếu đang giữ giá cao nhất → chuyển cho người giá cao thứ nhì

### 3.4 Trả lời câu hỏi

* Thực hiện tại trang chi tiết sản phẩm

### 3.5 Quản lý hồ sơ

* Xem danh sách sản phẩm đang đăng & còn hạn
* Xem danh sách sản phẩm đã có người thắng
* Đánh giá người thắng (+1 / -1)
* Được phép huỷ giao dịch và tự động -1 người thắng

---

## 4. Quản trị viên - administrator

### 4.1 Quản lý danh mục category

* CRUD
* Không xoá danh mục đã có sản phẩm

### 4.2 Quản lý sản phẩm

* Gỡ bỏ sản phẩm

### 4.3 Quản lý người dùng

* CRUD
* Xem danh sách bidder xin nâng cấp
* Duyệt nâng cấp bidder ➠ seller

### 4.4 Admin Dashboard

* Biểu đồ thống kê: sản phẩm mới, doanh thu, người dùng mới, bidder nâng cấp,…
* Các thống kê bổ sung tuỳ chọn

---

## 5. Các tính năng chung

### 5.1 Đăng nhập

* JWT Access/Refresh Token
* Có thể hỗ trợ Google, Facebook, Twitter, Github

### 5.2 Cập nhật thông tin cá nhân

* Họ tên, email, ngày sinh

### 5.3 Đổi mật khẩu

* bcrypt hoặc scrypt

### 5.4 Quên mật khẩu

* OTP qua email

---

## 6. Hệ thống

### 6.1 Mailing System

* Gửi email cho các sự kiện quan trọng:

  * Ra giá thành công
  * Người bị từ chối
  * Đấu giá kết thúc
  * Người mua đặt câu hỏi / người bán trả lời

### 6.2 Đấu giá tự động

* Người mua đặt giá-tối-đa
* Hệ thống tự nâng giá vừa đủ để thắng
* Nếu 2 người đặt cùng giá → ai đặt trước thắng
* Chỉ triển khai 1 trong 2: tự động hoặc thường

---

## 7. Quy trình thanh toán sau đấu giá

* 4 bước:

  1. Người mua thanh toán
  2. Người mua gửi địa chỉ
  3. Người bán xác nhận tiền và gửi hoá đơn vận chuyển
  4. Người mua xác nhận đã nhận hàng
* Hai bên có giao diện chat
* Cả hai có thể thay đổi đánh giá (+/-)
* Người bán có thể huỷ giao dịch và đánh giá -1

---

## 8. Các yêu cầu khác

### 8.1 Yêu cầu kỹ thuật

* Web App CSR
* Backend RESTful API + Swagger
* Logs, monitor (Grafana/ELK)
* Security: JWT
* Frontend SPA + router + validation + state management

### 8.2 Yêu cầu dữ liệu

* 20 sản phẩm, 4-5 danh mục
* Mỗi sản phẩm có ít nhất 5 lượt ra giá

### 8.3 Yêu cầu quản lý mã nguồn

* Upload code lên Github từ đầu
* Nhóm không có lịch sử commit → 0 điểm
