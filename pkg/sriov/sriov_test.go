package sriov

import (
	"errors"
	"github.com/Mellanox/ib-sriov-cni/pkg/types"
	"github.com/Mellanox/ib-sriov-cni/pkg/types/mocks"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/containernetworking/plugins/pkg/testutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/vishvananda/netlink"
)

// FakeLink is a dummy netlink struct used during testing
type FakeLink struct {
	netlink.LinkAttrs
}

// type FakeLink struct {
// 	linkAtrrs *netlink.LinkAttrs
// }

func (l *FakeLink) Attrs() *netlink.LinkAttrs {
	return &l.LinkAttrs
}

func (l *FakeLink) Type() string {
	return "FakeLink"
}

var _ = Describe("Sriov", func() {

	Context("Checking SetupVF function", func() {
		var (
			podifName string
			contID    string
			netconf   *types.NetConf
		)

		BeforeEach(func() {
			podifName = "net1"
			contID = "dummycid"
			netconf = &types.NetConf{
				Master:      "ib0",
				DeviceID:    "0000:af:06.0",
				VFID:        0,
				HostIFNames: "ib1",
				ContIFNames: "net1",
			}
		})

		It("Assuming existing interface", func() {
			var targetNetNS ns.NetNS
			targetNetNS, err := testutils.NewNS()
			defer func() {
				if targetNetNS != nil {
					targetNetNS.Close()
				}
			}()
			Expect(err).NotTo(HaveOccurred())
			mocked := &mocks.NetlinkManager{}

			fakeLink := &FakeLink{netlink.LinkAttrs{
				Index: 1000,
				Name:  "dummylink",
			}}

			mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
			mocked.On("LinkSetDown", fakeLink).Return(nil)
			mocked.On("LinkSetName", fakeLink, mock.Anything).Return(nil)
			mocked.On("LinkSetNsFd", fakeLink, mock.AnythingOfType("int")).Return(nil)
			mocked.On("LinkSetUp", fakeLink).Return(nil)
			sm := sriovManager{nLink: mocked}
			err = sm.SetupVF(netconf, podifName, contID, targetNetNS)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Assuming non existing interface", func() {
			var targetNetNS ns.NetNS
			targetNetNS, err := testutils.NewNS()
			defer func() {
				if targetNetNS != nil {
					targetNetNS.Close()
				}
			}()
			Expect(err).NotTo(HaveOccurred())
			mocked := &mocks.NetlinkManager{}

			mocked.On("LinkByName", mock.AnythingOfType("string")).Return(nil, errors.New("not fount"))
			sm := sriovManager{nLink: mocked}
			err = sm.SetupVF(netconf, podifName, contID, targetNetNS)
			Expect(err).To(HaveOccurred())
		})
		It("Assuming existing interface not able to set down", func() {
			var targetNetNS ns.NetNS
			targetNetNS, err := testutils.NewNS()
			defer func() {
				if targetNetNS != nil {
					targetNetNS.Close()
				}
			}()
			Expect(err).NotTo(HaveOccurred())
			mocked := &mocks.NetlinkManager{}

			fakeLink := &FakeLink{netlink.LinkAttrs{
				Index: 1000,
				Name:  "dummylink",
			}}

			mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
			mocked.On("LinkSetDown", fakeLink).Return(errors.New("failed"))
			sm := sriovManager{nLink: mocked}
			err = sm.SetupVF(netconf, podifName, contID, targetNetNS)
			Expect(err).To(HaveOccurred())
		})
		It("Assuming failed to change name", func() {
			var targetNetNS ns.NetNS
			targetNetNS, err := testutils.NewNS()
			defer func() {
				if targetNetNS != nil {
					targetNetNS.Close()
				}
			}()
			Expect(err).NotTo(HaveOccurred())
			mocked := &mocks.NetlinkManager{}

			fakeLink := &FakeLink{netlink.LinkAttrs{
				Index: 1000,
				Name:  "dummylink",
			}}

			mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
			mocked.On("LinkSetDown", fakeLink).Return(nil)
			mocked.On("LinkSetName", fakeLink, mock.Anything).Return(errors.New("failed"))
			sm := sriovManager{nLink: mocked}
			err = sm.SetupVF(netconf, podifName, contID, targetNetNS)
			Expect(err).To(HaveOccurred())
		})
		It("Assuming failed to move interface", func() {
			var targetNetNS ns.NetNS
			targetNetNS, err := testutils.NewNS()
			defer func() {
				if targetNetNS != nil {
					targetNetNS.Close()
				}
			}()
			Expect(err).ToNot(HaveOccurred())
			mocked := &mocks.NetlinkManager{}

			fakeLink := &FakeLink{netlink.LinkAttrs{
				Index: 1000,
				Name:  "dummylink",
			}}

			mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
			mocked.On("LinkSetDown", fakeLink).Return(nil)
			mocked.On("LinkSetName", fakeLink, mock.Anything).Return(nil)
			mocked.On("LinkSetNsFd", fakeLink, mock.AnythingOfType("int")).Return(errors.New("failed"))
			sm := sriovManager{nLink: mocked}
			err = sm.SetupVF(netconf, podifName, contID, targetNetNS)
			Expect(err).To(HaveOccurred())
		})
		It("Assuming failed to set interface up after moving", func() {
			var targetNetNS ns.NetNS
			targetNetNS, err := testutils.NewNS()
			defer func() {
				if targetNetNS != nil {
					targetNetNS.Close()
				}
			}()
			Expect(err).NotTo(HaveOccurred())
			mocked := &mocks.NetlinkManager{}

			fakeLink := &FakeLink{netlink.LinkAttrs{
				Index: 1000,
				Name:  "dummylink",
			}}

			mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
			mocked.On("LinkSetDown", fakeLink).Return(nil)
			mocked.On("LinkSetName", fakeLink, mock.Anything).Return(nil)
			mocked.On("LinkSetNsFd", fakeLink, mock.AnythingOfType("int")).Return(nil)
			mocked.On("LinkSetUp", fakeLink).Return(errors.New("failed"))
			sm := sriovManager{nLink: mocked}
			err = sm.SetupVF(netconf, podifName, contID, targetNetNS)
			Expect(err).To(HaveOccurred())
		})
		Context("Checking ReleaseVF function", func() {
			var (
				podifName string
				contID    string
				netconf   *types.NetConf
			)

			BeforeEach(func() {
				podifName = "net1"
				contID = "dummycid"
				netconf = &types.NetConf{
					Master:      "ib0",
					DeviceID:    "0000:af:06.0",
					VFID:        0,
					HostIFNames: "ib1",
					ContIFNames: "net1",
				}
			})
			It("Assuming existing interface", func() {
				var targetNetNS ns.NetNS
				targetNetNS, err := testutils.NewNS()
				defer func() {
					if targetNetNS != nil {
						targetNetNS.Close()
					}
				}()
				Expect(err).NotTo(HaveOccurred())
				mocked := &mocks.NetlinkManager{}
				fakeLink := &FakeLink{netlink.LinkAttrs{Index: 1000, Name: "dummylink"}}

				mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
				mocked.On("LinkSetDown", fakeLink).Return(nil)
				mocked.On("LinkSetName", fakeLink, mock.Anything).Return(nil)
				mocked.On("LinkSetNsFd", fakeLink, mock.AnythingOfType("int")).Return(nil)
				sm := sriovManager{nLink: mocked}
				err = sm.ReleaseVF(netconf, podifName, contID, targetNetNS)
				Expect(err).NotTo(HaveOccurred())
			})
			It("Assuming non existing interface", func() {
				var targetNetNS ns.NetNS
				targetNetNS, err := testutils.NewNS()
				defer func() {
					if targetNetNS != nil {
						targetNetNS.Close()
					}
				}()
				Expect(err).NotTo(HaveOccurred())
				mocked := &mocks.NetlinkManager{}

				mocked.On("LinkByName", mock.AnythingOfType("string")).Return(nil, errors.New("not found"))
				sm := sriovManager{nLink: mocked}
				err = sm.ReleaseVF(netconf, podifName, contID, targetNetNS)
				Expect(err).To(HaveOccurred())
			})
			It("Assuming failed to set interface down", func() {
				var targetNetNS ns.NetNS
				targetNetNS, err := testutils.NewNS()
				defer func() {
					if targetNetNS != nil {
						targetNetNS.Close()
					}
				}()
				Expect(err).NotTo(HaveOccurred())
				mocked := &mocks.NetlinkManager{}
				fakeLink := &FakeLink{netlink.LinkAttrs{Index: 1000, Name: "dummylink"}}

				mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
				mocked.On("LinkSetDown", fakeLink).Return(errors.New("failed"))
				sm := sriovManager{nLink: mocked}
				err = sm.ReleaseVF(netconf, podifName, contID, targetNetNS)
				Expect(err).To(HaveOccurred())
			})
			It("Assuming failed to move interface", func() {
				var targetNetNS ns.NetNS
				targetNetNS, err := testutils.NewNS()
				defer func() {
					if targetNetNS != nil {
						targetNetNS.Close()
					}
				}()
				Expect(err).NotTo(HaveOccurred())
				mocked := &mocks.NetlinkManager{}
				fakeLink := &FakeLink{netlink.LinkAttrs{Index: 1000, Name: "dummylink"}}

				mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
				mocked.On("LinkSetDown", fakeLink).Return(nil)
				mocked.On("LinkSetName", fakeLink, mock.Anything).Return(errors.New("failed"))
				sm := sriovManager{nLink: mocked}
				err = sm.ReleaseVF(netconf, podifName, contID, targetNetNS)
				Expect(err).To(HaveOccurred())
			})
			It("Assuming existing interface", func() {
				var targetNetNS ns.NetNS
				targetNetNS, err := testutils.NewNS()
				defer func() {
					if targetNetNS != nil {
						targetNetNS.Close()
					}
				}()
				Expect(err).NotTo(HaveOccurred())
				mocked := &mocks.NetlinkManager{}
				fakeLink := &FakeLink{netlink.LinkAttrs{Index: 1000, Name: "dummylink"}}

				mocked.On("LinkByName", mock.AnythingOfType("string")).Return(fakeLink, nil)
				mocked.On("LinkSetDown", fakeLink).Return(nil)
				mocked.On("LinkSetName", fakeLink, mock.Anything).Return(nil)
				mocked.On("LinkSetNsFd", fakeLink, mock.AnythingOfType("int")).Return(errors.New("failed"))
				sm := sriovManager{nLink: mocked}
				err = sm.ReleaseVF(netconf, podifName, contID, targetNetNS)
				Expect(err).To(HaveOccurred())
			})

		})
	})
})
