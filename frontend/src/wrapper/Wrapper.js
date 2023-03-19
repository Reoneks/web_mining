import { BrowserRouter, Routes, Route } from "react-router-dom";
import WrappedComponent from "../parsed/App";
import Styles from "./styles";
import { Form, Field } from "react-final-form";
import arrayMutators from "final-form-arrays";
import { FieldArray } from "react-final-form-arrays";
import { Link } from "react-router-dom";

export default function Wrapper() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={
            <Styles>
              <h1>Parse site</h1>
              <Form
                onSubmit={function () {}}
                mutators={{
                  ...arrayMutators,
                }}
                render={({
                  handleSubmit,
                  form: {
                    mutators: { push, pop },
                  }, // injected from final-form-arrays above
                  pristine,
                  form,
                  submitting,
                  values,
                }) => {
                  return (
                    <form onSubmit={handleSubmit}>
                      <div>
                        <label>Site url</label>
                        <Field name="url" component="input" />
                      </div>
                      <div>
                        <label>Only this page</label>
                        <Field
                          name="only_this_page"
                          component="input"
                          type="checkbox"
                        />
                      </div>
                      <div>
                        <label>Force collect</label>
                        <Field name="force" component="input" type="checkbox" />
                      </div>
                      <div className="buttons">
                        <button
                          type="button"
                          onClick={() => push("exclude", undefined)}
                        >
                          Exclude url
                        </button>
                        <button type="button" onClick={() => pop("exclude")}>
                          Remove exclude
                        </button>
                      </div>
                      <FieldArray name="exclude">
                        {({ fields }) =>
                          fields.map((url, index) => (
                            <div key={url}>
                              <label>Url №{index + 1}</label>
                              <Field
                                name={`${url}`}
                                component="input"
                                placeholder="URL"
                              />
                              <span
                                onClick={() => fields.remove(index)}
                                style={{ cursor: "pointer" }}
                              >
                                ❌
                              </span>
                            </div>
                          ))
                        }
                      </FieldArray>

                      <div className="buttons">
                        <Link to="/parsed" state={values}>
                          <button
                            type="submit"
                            disabled={submitting || pristine}
                          >
                            Submit
                          </button>
                        </Link>
                        <button
                          type="button"
                          onClick={form.reset}
                          disabled={submitting || pristine}
                        >
                          Reset
                        </button>
                      </div>
                    </form>
                  );
                }}
              />
            </Styles>
          }
        />
        <Route path="/parsed" element={<WrappedComponent />} />
      </Routes>
    </BrowserRouter>
  );
}
