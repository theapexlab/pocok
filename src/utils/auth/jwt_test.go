package auth_test

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pocok/src/utils/auth"
	"pocok/src/utils/models"
)

// Override time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

// this is a parsed token, created at unix time 0, valid for 2 days containing testOrg
var testToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjgwMCwib3JnSWQiOiJURVNUX09SR0FOSVpBVElPTiJ9.vMVdagL7KXcmudO4O1M6_pvEvHC2uLfyoH9QtTqoOJU"
var testOrg string = "TEST_ORGANIZATION"
var testJwtKey string = "TEST_KICSI_CICA"

var _ = Describe("Auth", func() {
	var err error
	var token string
	var payload *models.JWTClaims
	os.Setenv("jwtKey", testJwtKey)

	When("Creating token", func() {
		BeforeEach(func() {
			at(time.Unix(0, 0), func() {
				token, err = CreateToken(testOrg)
			})
		})

		It("should not error", func() {
			Expect(err).To(BeNil())
		})

		It("should give back encrypted token", func() {
			Expect(token).To(Equal(testToken))
		})
	})

	When("parsing valid token", func() {
		BeforeEach(func() {
			at(time.Unix(0, 0), func() {
				payload, err = ParseToken(testToken)
			})
		})

		It("should not error", func() {
			Expect(err).To(BeNil())
		})

		It("should give back encrypted token", func() {
			Expect(payload.OrgId).To(Equal(testOrg))
		})
	})

	When("parsing expired token", func() {
		BeforeEach(func() {
			at(time.Unix(172801, 0), func() {
				payload, err = ParseToken(testToken)
			})
		})

		It("should have error", func() {
			Expect(err).To(MatchError("token is expired by 1s"))
		})

		It("should not give back payload", func() {
			Expect(payload).To(BeNil())
		})
	})

	When("creating and parsing from token", func() {
		var createErr error
		var parseErr error
		testPayload := time.Now().GoString()
		BeforeEach(func() {
			token, createErr = CreateToken(testPayload)
			payload, parseErr = ParseToken(token)
		})

		It("should not have any error", func() {
			Expect(createErr).To(BeNil())
			Expect(parseErr).To(BeNil())
		})

		It("should not give back payload", func() {
			Expect(payload.OrgId).To(Equal(testPayload))
		})
	})

})
