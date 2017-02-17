package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

var expect []Target

func setup() {
	fmt.Println("setup")

}

func tearDown() {
	fmt.Println("tearDown")
}

type testDate struct {
	path   string
	expect []Target
}

func TestParseFileForAppBuildGradle(t *testing.T) {

	testdate := []testDate{
		testDate{
			"testdata/app/build.gradle",
			[]Target{
				Target{"junit", "junit", "4.12"},
				Target{"com.android.support", "appcompat-v7", "${support_lib_version}"},
				Target{"com.android.support", "recyclerview-v7", "${support_lib_version}"},
				Target{"com.android.support", "design", "${support_lib_version}"},
				Target{"com.android.support", "support-v4", "${support_lib_version}"},
				Target{"com.google.android", "flexbox", "0.2.3"},
				Target{"com.google.code.gson", "gson", "2.7"},
				Target{"com.squareup.okhttp3", "okhttp", "${okhttp_version}"},
				Target{"com.squareup.okhttp3", "logging-interceptor", "${okhttp_version}"},
				Target{"io.reactivex", "rxjava", "1.2.6"},
				Target{"io.reactivex", "rxandroid", "1.2.1"},
				Target{"com.squareup.retrofit2", "retrofit", "${retrofit_version}"},
				Target{"com.squareup.retrofit2", "converter-gson", "${retrofit_version}"},
				Target{"com.squareup.retrofit2", "adapter-rxjava", "${retrofit_version}"},
				Target{"com.google.firebase", "firebase-core", "${google_play_services_version}"},
				Target{"com.google.firebase", "firebase-ads", "${google_play_services_version}"},
				Target{"com.google.firebase", "firebase-crash", "${google_play_services_version}"},
				Target{"com.google.firebase", "firebase-auth", "${google_play_services_version}"},
				Target{"com.google.android.gms", "play-services-analytics", "${google_play_services_version}"},
				Target{"com.jakewharton.rxbinding", "rxbinding-recyclerview-v7", "1.0.0"},
				Target{"com.annimon", "stream", "1.1.5"},
				Target{"com.github.bumptech.glide", "glide", "3.7.0"},
				Target{"com.roughike", "bottom-bar", "2.0.2"},
				Target{"com.google.dagger", "dagger", "${dagger_version}"},
				Target{"com.android.support.test", "runner", "0.5"},
				Target{"com.android.support.test", "rules", "0.5"},
				Target{"com.squareup.okhttp3", "mockwebserver", "${okhttp_version}"},
				Target{"com.android.support", "support-annotations", "${support_lib_version}"},
				Target{"com.squareup.leakcanary", "leakcanary-android", "${leakcanary_version}"},
				Target{"com.squareup.leakcanary", "leakcanary-android-no-op", "${leakcanary_version}"},
				Target{"com.squareup.leakcanary", "leakcanary-android-no-op", "${leakcanary_version}"},
			},
		},
		testDate{
			"testdata/build.gradle",
			[]Target{
				Target{"com.android.tools.build", "gradle", "2.2.3"},
				Target{"com.google.gms", "google-services", "3.0.0"},
				Target{"com.github.triplet.gradle", "play-publisher", "1.1.5"},
				Target{"me.tatarka", "gradle-retrolambda", "3.5.0"},
			},
		},
	}

	for _, td := range testdate {
		deps, err := parseFile(td.path)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(td.expect, deps) {
			t.Errorf("expect:%v, \ndeps:%v", td.expect, deps)
		}

	}

}
