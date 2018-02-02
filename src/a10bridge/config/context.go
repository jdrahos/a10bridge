package config

type RunContext struct {
	Arguments    *Args
	A10Instances A10Instances
}

func BuildConfig() (*RunContext, error) {
	var context *RunContext
	args, err := buildArguments()

	if err != nil {
		return context, err
	}

	a10Config, err := readA10Configuration(*args.A10Config)
	if err != nil {
		return context, err
	}

	instances := make([]A10Instance, 0)

	for _, instance := range a10Config.Instances {
		if len(instance.Name) == 0 {
			instance.Name = instance.APIUrl
		}
		if len(instance.Password) == 0 {
			instance.Password = *args.A10Pwd
		}
		instances = append(instances, instance)
	}

	return &RunContext{
		Arguments:    args,
		A10Instances: instances,
	}, err
}
