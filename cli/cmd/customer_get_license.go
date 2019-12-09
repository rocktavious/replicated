package cmd

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

func (r *runners) InitCustomerGetLicense(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-license CUSTOMER",
		Short: "Get a license for a customer",
		Long:  `Get a license for a customer. Use with -o FILENAME `,
		RunE:  r.createCustomer,
	}
	parent.AddCommand(cmd)
	cmd.Flags().StringVarP(&r.args.customerGetLicenseOutputPath, "output", "o", "", "Output Path. Default is $CUSTOMER-NAME.yam")

	return cmd
}

func (r *runners) getLicense(_ *cobra.Command, args []string) error {

	if len(args) != 1 {
		return errors.New("Exactly one customer name or ID required")
	}

	customerNameOrID := args[0]
	customer, err := r.api.GetCustomerByNameOrID(r.appID, r.appType, customerNameOrID)
	if err != nil {
		return errors.Wrap(err, "find customer")
	}

	license, err := r.api.GetCustomerLicense(r.appID, r.appType, customer.ID)


	filename := r.args.customerGetLicenseOutputPath
	if filename == "-" {
		_, err = bytes.NewReader(license).WriteTo(os.Stdout)
		return errors.Wrap(err, "write to stdout")
	}


	if filename == "" {
		filename = fmt.Sprintf("%s-%s.yaml", customer.Name, r.appName)
	}

	err = ioutil.WriteFile(filename, license, 0644)
	if err != nil {
		return errors.Wrap(err, "write license")
	}

	_, err = fmt.Fprintf(os.Stderr, "License file written to %s", filename)

	return err
}
