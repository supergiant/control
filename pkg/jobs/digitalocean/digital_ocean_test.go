package digitalocean

import (
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/mock"

	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/testutils"
)

type createDropletResponse struct {
	droplet *godo.Droplet
	resp    *godo.Response
	err     error
}

type getDropletResponse struct {
	droplet *godo.Droplet
	resp    *godo.Response
	err     error
}

type (
	fakeDropletService struct {
		createCounter int
		getCounter    int

		createResponses []createDropletResponse
		getResponses    []getDropletResponse
	}
	fakeTagService struct {
		counter int
		errs    []error
	}
)

func (f *fakeDropletService) Get(id int) (*godo.Droplet, *godo.Response, error) {
	if f.getCounter > len(f.getResponses)-1 {
		panic("get index out of range")
	}

	r := f.getResponses[f.getCounter]
	f.getCounter++
	return r.droplet, r.resp, r.err
}

func (f *fakeDropletService) Create(req *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
	if f.createCounter > len(f.createResponses)-1 {
		panic("create index out of range")
	}
	r := f.createResponses[f.createCounter]
	f.createCounter++
	return r.droplet, r.resp, r.err
}

func (f *fakeTagService) TagResources(string, *godo.TagResourcesRequest) (*godo.Response, error) {
	if f.counter > len(f.errs)-1 {
		panic("tag index out of range")
	}

	return nil, f.errs[f.counter]
}

func TestJob_CreateDropletSuccess(t *testing.T) {
	dropletId := 1
	m := new(testutils.MockStorage)
	m.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	errTaggingDroplet := errors.New("error tagging droplet")

	testCases := []struct {
		storage        storage.Interface
		dropletService *fakeDropletService
		tagService     *fakeTagService

		dropletTimeout time.Duration
		checkTimeout   time.Duration

		expectedError error
	}{
		{
			storage: m,
			dropletService: &fakeDropletService{
				0,
				0,
				[]createDropletResponse{
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
				},
				[]getDropletResponse{
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "active",
						},
						nil,
						nil,
					},
				},
			},
			tagService: &fakeTagService{
				0,
				[]error{nil},
			},
			dropletTimeout: time.Millisecond * 30,
			checkTimeout:   time.Millisecond * 1,
			expectedError:  nil,
		},
		{
			storage: m,
			dropletService: &fakeDropletService{
				0,
				0,
				[]createDropletResponse{
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
				},
				[]getDropletResponse{
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
				},
			},
			tagService: &fakeTagService{
				0,
				[]error{nil},
			},
			dropletTimeout: time.Millisecond * 2,
			checkTimeout:   time.Millisecond * 1,
			expectedError:  ErrTimeoutExceeded,
		},
		{
			storage: m,
			dropletService: &fakeDropletService{
				0,
				0,
				[]createDropletResponse{
					{
						&godo.Droplet{
							ID:     dropletId,
							Name:   "test",
							Status: "pending",
						},
						nil,
						nil,
					},
				},
				[]getDropletResponse{},
			},
			tagService: &fakeTagService{
				0,
				[]error{
					errTaggingDroplet,
				},
			},
			dropletTimeout: time.Millisecond * 2,
			checkTimeout:   time.Millisecond * 1,
			expectedError:  errTaggingDroplet,
		},
	}

	for _, testCase := range testCases {
		job := &Job{
			storage:        testCase.storage,
			dropletService: testCase.dropletService,
			tagService:     testCase.tagService,
		}

		config := Config{
			"test",
			"1.8.7",
			"us-west1",
			"2GB",
			"master",
			[]string{"fingerprint"},

			testCase.dropletTimeout,
			testCase.checkTimeout,
		}

		err := job.CreateDroplet(config)

		if err != testCase.expectedError {
			t.Errorf("Wrong error expected %s actual %s", testCase.expectedError.Error(), err.Error())
		}
	}
}
