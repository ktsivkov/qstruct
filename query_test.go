package qstruct_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ktsivkov/qstruct"
)

func TestNewFor__invalidType(t *testing.T) {
	type demoStruct struct {
		Value int `query:"value" validate:"required"`
	}

	values, _ := url.ParseQuery("")
	_, err := qstruct.NewFor[[]demoStruct](values)
	assert.ErrorIs(t, err, qstruct.ErrUnexpectedType)
}

func TestNewFor__validate(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		t.Run("fail", func(t *testing.T) {
			type demoStruct struct {
				Value int `query:"value" validate:"required"`
			}

			values, _ := url.ParseQuery("")
			res, err := qstruct.NewFor[demoStruct](values)
			assert.Nil(t, res)
			assert.ErrorIs(t, err, qstruct.ErrRequired)
		})
		t.Run("success", func(t *testing.T) {
			type demoStruct struct {
				Value int `query:"value" validate:"required"`
			}

			values, _ := url.ParseQuery("value=12")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.NoError(t, err)
		})
	})
}

func TestNewFor__ignore(t *testing.T) {
	type demoStruct struct {
		Value int `query:"-"`
	}

	values, _ := url.ParseQuery("value=12")
	res, err := qstruct.NewFor[demoStruct](values)
	assert.NoError(t, err)
	assert.Equal(t, &demoStruct{}, res)
}

func TestNewFor__unnamed(t *testing.T) {
	type demoStruct struct {
		Value int
	}

	values, _ := url.ParseQuery("Value=12")
	res, err := qstruct.NewFor[demoStruct](values)
	assert.NoError(t, err)
	assert.Equal(t, &demoStruct{Value: 12}, res)
}

func TestNewFor__int(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value int `query:"value"`
			}

			values, _ := url.ParseQuery("value=12")
			res, err := qstruct.NewFor[demoStruct](values)
			assert.NoError(t, err)
			assert.Equal(t, &demoStruct{Value: 12}, res)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value int `query:"value" default:"12"`
			}

			values, _ := url.ParseQuery("")
			res, err := qstruct.NewFor[demoStruct](values)
			assert.NoError(t, err)
			assert.Equal(t, &demoStruct{Value: 12}, res)
		})
	})
	t.Run("invalid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value int `query:"value"`
			}

			values, _ := url.ParseQuery("value=abc")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value int `query:"value" default:"abc"`
			}

			values, _ := url.ParseQuery("")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
	})
}

func TestNewFor__uint(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value uint `query:"value"`
			}

			values, _ := url.ParseQuery("value=12")
			res, err := qstruct.NewFor[demoStruct](values)
			assert.NoError(t, err)
			assert.Equal(t, &demoStruct{Value: 12}, res)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value uint `query:"value" default:"12"`
			}

			values, _ := url.ParseQuery("")
			res, err := qstruct.NewFor[demoStruct](values)
			assert.NoError(t, err)
			assert.Equal(t, &demoStruct{Value: 12}, res)
		})
	})
	t.Run("invalid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value uint `query:"value"`
			}

			values, _ := url.ParseQuery("value=-12")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value uint `query:"value" default:"-12"`
			}

			values, _ := url.ParseQuery("")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
	})
}

func TestNewFor__string(t *testing.T) {
	t.Run("from given", func(t *testing.T) {
		type demoStruct struct {
			Value string `query:"value"`
		}

		values, _ := url.ParseQuery("value=my value")
		res, err := qstruct.NewFor[demoStruct](values)
		assert.NoError(t, err)
		assert.Equal(t, &demoStruct{Value: "my value"}, res)
	})
	t.Run("from default", func(t *testing.T) {
		type demoStruct struct {
			Value string `query:"value" default:"my value"`
		}

		values, _ := url.ParseQuery("")
		res, err := qstruct.NewFor[demoStruct](values)
		assert.NoError(t, err)
		assert.Equal(t, &demoStruct{Value: "my value"}, res)
	})
}

func TestNewFor__float(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value float32 `query:"value"`
			}

			values, _ := url.ParseQuery("value=12.12")
			res, err := qstruct.NewFor[demoStruct](values)
			assert.NoError(t, err)
			assert.Equal(t, &demoStruct{Value: 12.12}, res)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value float32 `query:"value" default:"12.12"`
			}

			values, _ := url.ParseQuery("")
			res, err := qstruct.NewFor[demoStruct](values)
			assert.NoError(t, err)
			assert.Equal(t, &demoStruct{Value: 12.12}, res)
		})
	})
	t.Run("invalid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value float32 `query:"value"`
			}

			values, _ := url.ParseQuery("value=abc")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value float32 `query:"value" default:"abc"`
			}

			values, _ := url.ParseQuery("")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
	})
}

func TestNewFor__bool(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("from bool", func(t *testing.T) {
			t.Run("for true", func(t *testing.T) {
				t.Run("from given", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value"`
					}

					values, _ := url.ParseQuery("value=true")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: true}, res)
				})
				t.Run("from default", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value" default:"true"`
					}

					values, _ := url.ParseQuery("")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: true}, res)
				})
			})
			t.Run("for false", func(t *testing.T) {
				t.Run("from given", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value"`
					}

					values, _ := url.ParseQuery("value=false")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: false}, res)
				})
				t.Run("from default", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value" default:"false"`
					}

					values, _ := url.ParseQuery("")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: false}, res)
				})
			})
		})
		t.Run("from literal", func(t *testing.T) {
			t.Run("for true", func(t *testing.T) {
				t.Run("from given", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value"`
					}

					values, _ := url.ParseQuery("value=1")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: true}, res)
				})
				t.Run("from default", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value" default:"1"`
					}

					values, _ := url.ParseQuery("")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: true}, res)
				})
			})
			t.Run("for false", func(t *testing.T) {
				t.Run("from given", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value"`
					}

					values, _ := url.ParseQuery("value=0")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: false}, res)
				})
				t.Run("from default", func(t *testing.T) {
					type demoStruct struct {
						Value bool `query:"value" default:"0"`
					}

					values, _ := url.ParseQuery("")
					res, err := qstruct.NewFor[demoStruct](values)
					assert.NoError(t, err)
					assert.Equal(t, &demoStruct{Value: false}, res)
				})
			})
		})
	})
	t.Run("invalid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value bool `query:"value"`
			}

			values, _ := url.ParseQuery("value=abc")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value bool `query:"value" default:"abc"`
			}

			values, _ := url.ParseQuery("")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
	})
}

func TestNewFor__time(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("default format", func(t *testing.T) {
			t.Run("from given", func(t *testing.T) {
				type demoStruct struct {
					Value time.Time `query:"value"`
				}

				values, _ := url.ParseQuery("value=2000-12-31T00:00:00Z")
				res, err := qstruct.NewFor[demoStruct](values)
				assert.NoError(t, err)
				expected, _ := time.Parse(time.RFC3339, "2000-12-31T00:00:00Z")
				assert.Equal(t, &demoStruct{Value: expected}, res)
			})
			t.Run("from default", func(t *testing.T) {
				type demoStruct struct {
					Value time.Time `query:"value" default:"2000-12-31T00:00:00Z"`
				}

				values, _ := url.ParseQuery("")
				res, err := qstruct.NewFor[demoStruct](values)
				assert.NoError(t, err)
				expected, _ := time.Parse(time.RFC3339, "2000-12-31T00:00:00Z")
				assert.Equal(t, &demoStruct{Value: expected}, res)
			})
		})
		t.Run("custom format", func(t *testing.T) {
			t.Run("from given", func(t *testing.T) {
				type demoStruct struct {
					Value time.Time `query:"value" format:"Mon, 02 Jan 2006 15:04:05 -0700"`
				}

				values, _ := url.ParseQuery("value=Sun, 31 Dec 2000 00:00:00 %2B0000")
				res, err := qstruct.NewFor[demoStruct](values)
				assert.NoError(t, err)
				expected, _ := time.Parse(time.RFC1123Z, "Sun, 31 Dec 2000 00:00:00 +0000")
				assert.Equal(t, &demoStruct{Value: expected}, res)
			})
			t.Run("from default", func(t *testing.T) {
				type demoStruct struct {
					Value time.Time `query:"value" default:"Sun, 31 Dec 2000 00:00:00 +0000" format:"Mon, 02 Jan 2006 15:04:05 -0700"`
				}

				values, _ := url.ParseQuery("")
				res, err := qstruct.NewFor[demoStruct](values)
				assert.NoError(t, err)
				expected, _ := time.Parse(time.RFC1123Z, "Sun, 31 Dec 2000 00:00:00 +0000")
				assert.Equal(t, &demoStruct{Value: expected}, res)
			})
		})
	})
	t.Run("invalid", func(t *testing.T) {
		t.Run("from given", func(t *testing.T) {
			type demoStruct struct {
				Value time.Time `query:"value"`
			}

			values, _ := url.ParseQuery("value=abc")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
		t.Run("from default", func(t *testing.T) {
			type demoStruct struct {
				Value time.Time `query:"value" default:"abc"`
			}

			values, _ := url.ParseQuery("")
			_, err := qstruct.NewFor[demoStruct](values)
			assert.ErrorIs(t, err, qstruct.ErrUnexpectedValue)
		})
	})
}
