[![wercker status](https://app.wercker.com/status/6d9e484f12bbb5152302b39e02593349/s/master "wercker status")](https://app.wercker.com/project/byKey/6d9e484f12bbb5152302b39e02593349)
[![Go Report Card](https://goreportcard.com/badge/github.com/kwmt/libver)](https://goreportcard.com/report/github.com/kwmt/libver)


This is a tool that fetch latest library version in build.gradle file of Android.


## Usage 

You need to set two environment variables.

* `BINTRAY_API_USERNAME`
* `BINTRAY_API_PASSWORD`

user’s name of [bintray](https://bintray.com)  as username and the API key as the password. You can obtain your API key from [bintray profile page](https://bintray.com/profile/edit).

In your Android project top directory, 

```
$ go get github.com/kwmt/libver
$ libver app/build.gradle
```

for [app/build](https://github.com/kwmt/libver/blob/master/testdata/app/build.gradle), output like below:
```
junit:junit:4.12-beta-3
com.android.support:25.1.1
com.google.code.gson:gson:2.8.0
com.squareup.okhttp3:okhttp:3.6.0
com.squareup.okhttp3:logging-interceptor:3.6.0
RxJava:2.0.6
RxAndroid:2.0.1
com.squareup.retrofit2:retrofit:2.1.0
com.squareup.retrofit2:converter-gson:2.1.0
com.squareup.retrofit2:adapter-rxjava:2.1.0
com.jakewharton.rxbinding:0.3.2.8
com.annimon:stream:1.1.5
glide-library:3.4.0.1
bottom-bar:2.1.1
com.google.dagger:dagger:2.9
com.squareup.okhttp3:mockwebserver:3.6.0
com.squareup.leakcanary:leakcanary-android:1.3
com.squareup.leakcanary:leakcanary-android-no-op:1.3
com.squareup.leakcanary:leakcanary-android-no-op:1.3
```

## With CI

Offcourse, you can use on CI service.

For example on Wercker CI below:

<img width="558" alt="2017-02-17 21 43 13" src="https://cloud.githubusercontent.com/assets/1450486/23065688/3282fe64-f55a-11e6-9ee9-ed76df8e3e62.png">

### wercker sample

https://github.com/kwmt/GitHubSearch/blob/master/wercker.yml#L14..L33


```
    - script:
        name: install git
        code: |
          sudo apt-get -y install git
    - script:
        name: install golang
        code: |
          wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
          sudo tar -xvf go1.8.linux-amd64.tar.gz
          sudo mv go /usr/local
          export GOROOT=/usr/local/go
          export GOPATH=$PWD
          export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
          go version
    - script:
        name: check dependencies library version
        code: |
          go get github.com/kwmt/libver
          libver app/build.gradle
```

## TODO

- [ ] fetch library include Android SDK(ex. "com.google.android", "com.google.firebase" etc...)
- [ ] latest version may not latest version. (specification of bintray?) so **please just for reference**

