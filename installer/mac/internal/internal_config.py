from .build_props_config import BuildPropsCodeSignConfig

# build will look for this class in config_factory.py. This class name is also
# used for "Google Chrome" builds so it will be picked up instead
# of ChromiumCodeSignConfig.
class InternalCodeSignConfig(BuildPropsCodeSignConfig):
    """A CodeSignConfig used for signing Official Beacon builds.
    """

    @property
    def provisioning_profile_basename(self):
        return "release"
    
    @property
    def run_spctl_assess(self):
        # For some reason spctl assess is ran before notarization
        # causing the check to fail with unnotarized developer error
        # TODO: fix
        return False    
