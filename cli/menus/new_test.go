package menus

import (
	"testing"

	"github.com/razshare/frizzante/cli/apps"
)

func TestNew(t *testing.T) {
	a := apps.New()
	menu, err := New(a)
	if err != nil {
		t.Fatal(err)
	}

	i := 1

	for _, item := range menu.Items {
		if item.Choice.Id == "configure" {
			i++
			if item.Active() {
				t.Fatal("configure should not be active")
			}
			*a.Configure = true
			if !item.Active() {
				t.Fatal("configure should be active")
			}
			*a.Configure = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "create project" {
			i++
			if item.Active() {
				t.Fatal("create project should not be active")
			}
			*a.CreateProject = "asd"
			if !item.Active() {
				t.Fatal("create project should be active")
			}
			*a.CreateProject = ""
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "install" {
			i++
			if item.Active() {
				t.Fatal("install should not be active")
			}
			*a.Install = true
			if !item.Active() {
				t.Fatal("install should be active")
			}
			*a.Install = false
			continue
		}

		if item.Choice.Id == "install" {
			i++
			if item.Active() {
				t.Fatal("install should not be active")
			}
			*a.Install = true
			if !item.Active() {
				t.Fatal("install should be active")
			}
			*a.Install = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "update" {
			i++
			if item.Active() {
				t.Fatal("update should not be active")
			}
			*a.Update = true
			if !item.Active() {
				t.Fatal("update should be active")
			}
			*a.Update = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "add" {
			i++
			if item.Active() {
				t.Fatal("add should not be active")
			}
			*a.Add = "test"
			if !item.Active() {
				t.Fatal("add should be active")
			}
			*a.Add = ""
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "sveltekit-sse" {
			i++
			if item.Active() {
				t.Fatal("dev should not be active")
			}
			*a.Dev = true
			if !item.Active() {
				t.Fatal("dev should be active")
			}
			*a.Dev = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "build" {
			i++
			if item.Active() {
				t.Fatal("build should not be active")
			}
			*a.Build = true
			if !item.Active() {
				t.Fatal("build should be active")
			}
			*a.Build = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "generate" {
			i++
			if item.Active() {
				t.Fatal("generate should not be active")
			}
			*a.Generate = true
			if !item.Active() {
				t.Fatal("generate should be active")
			}
			*a.Generate = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "package" {
			i++
			if item.Active() {
				t.Fatal("package should not be active")
			}
			*a.Package = true
			if !item.Active() {
				t.Fatal("package should be active")
			}
			*a.Package = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "package (watch)" {
			i++
			if item.Active() {
				t.Fatal("package (watch) should not be active")
			}
			*a.PackageWatch = true
			if !item.Active() {
				t.Fatal("package (watch) should be active")
			}
			*a.PackageWatch = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "check" {
			i++
			if item.Active() {
				t.Fatal("check should not be active")
			}
			*a.Check = true
			if !item.Active() {
				t.Fatal("check should be active")
			}
			*a.Check = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "format" {
			i++
			if item.Active() {
				t.Fatal("format should not be active")
			}
			*a.Format = true
			if !item.Active() {
				t.Fatal("format should be active")
			}
			*a.Format = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "touch" {
			i++
			if item.Active() {
				t.Fatal("touch should not be active")
			}
			*a.Touch = true
			if !item.Active() {
				t.Fatal("touch should be active")
			}
			*a.Touch = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "clean project" {
			i++
			if item.Active() {
				t.Fatal("clean project should not be active")
			}
			*a.CleanProject = true
			if !item.Active() {
				t.Fatal("clean project should be active")
			}
			*a.CleanProject = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "reset" {
			i++
			if item.Active() {
				t.Fatal("reset should not be active")
			}
			*a.Reset = true
			if !item.Active() {
				t.Fatal("reset should be active")
			}
			*a.Reset = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "clear" {
			i++
			if item.Active() {
				t.Fatal("clear should not be active")
			}
			*a.Clear = true
			if !item.Active() {
				t.Fatal("clear should be active")
			}
			*a.Clear = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "test" {
			i++
			if item.Active() {
				t.Fatal("test should not be active")
			}
			*a.Test = true
			if !item.Active() {
				t.Fatal("test should be active")
			}
			*a.Test = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "welcome" {
			i++
			if item.Active() {
				t.Fatal("welcome should not be active")
			}
			*a.Welcome = true
			if !item.Active() {
				t.Fatal("welcome should be active")
			}
			*a.Welcome = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "help" {
			i++
			if item.Active() {
				t.Fatal("help should not be active")
			}
			*a.Help = true
			if !item.Active() {
				t.Fatal("help should be active")
			}
			*a.Help = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "version" {
			i++
			if item.Active() {
				t.Fatal("version should not be active")
			}
			*a.Version = true
			if !item.Active() {
				t.Fatal("version should be active")
			}
			*a.Version = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "assembly explorer" {
			i++
			if item.Active() {
				t.Fatal("assembly explorer should not be active")
			}
			*a.AssemblyExplorer = true
			if !item.Active() {
				t.Fatal("assembly explorer should be active")
			}
			*a.AssemblyExplorer = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}

		if item.Choice.Id == "migrate" {
			i++
			if item.Active() {
				t.Fatal("migrate should not be active")
			}
			*a.Migrate = true
			if !item.Active() {
				t.Fatal("migrate should be active")
			}
			*a.Migrate = false
			//err = item.Handler()
			//if err != nil {
			//	t.Fatal(err)
			//}
			continue
		}
	}

	if i != len(menu.Items) {
		t.Fatal("some menu items were skipped")
	}
}
