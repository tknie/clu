package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ogen-go/ogen/uri"
	"github.com/tknie/log"
)

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	elemIsEscaped := false
	if rawPath := r.URL.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
			elemIsEscaped = strings.ContainsRune(elem, '%')
		}
	}
	log.Log.Debugf("Router elem before prefix  %s", elem)
	log.Log.Debugf("Router prefix %s", s.cfg.Prefix)
	if prefix := s.cfg.Prefix; len(prefix) > 0 {
		if strings.HasPrefix(elem, prefix) {
			// Cut prefix from the path.
			elem = strings.TrimPrefix(elem, prefix)
		} else {
			log.Log.Debugf("Prefix doesn't match %s", elem)
			// Prefix doesn't match.
			s.notFound(w, r)
			return
		}
	}
	log.Log.Debugf("Router elem after prefix  %s", elem)
	if len(elem) == 0 {
		log.Log.Debugf("Elem empty %s", elem)
		s.notFound(w, r)
		return
	}
	args := [3]string{}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'a': // Prefix: "admin/access/"
				if l := len("admin/access/"); len(elem) >= l && elem[0:l] == "admin/access/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "role"
				// Leaf parameter
				args[0] = elem
				elem = ""

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "DELETE":
						s.handleDelAccessRequest([1]string{
							args[0],
						}, elemIsEscaped, w, r)
					case "GET":
						s.handleAccessRequest([1]string{
							args[0],
						}, elemIsEscaped, w, r)
					case "POST":
						s.handleAddAccessRequest([1]string{
							args[0],
						}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "DELETE,GET,POST")
					}

					return
				}
			case 'b': // Prefix: "binary/"
				if l := len("binary/"); len(elem) >= l && elem[0:l] == "binary/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "table"
				// Match until "/"
				idx := strings.IndexByte(elem, '/')
				if idx < 0 {
					idx = len(elem)
				}
				args[0] = elem[:idx]
				elem = elem[idx:]

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "field"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[1] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "search"
						// Leaf parameter
						args[2] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleGetLobByMapRequest([3]string{
									args[0],
									args[1],
									args[2],
								}, elemIsEscaped, w, r)
							case "PUT":
								s.handleUpdateLobByMapRequest([3]string{
									args[0],
									args[1],
									args[2],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET,PUT")
							}

							return
						}
					}
				}
			case 'c': // Prefix: "config"
				if l := len("config"); len(elem) >= l && elem[0:l] == "config" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch r.Method {
					case "GET":
						s.handleGetConfigRequest([0]string{}, elemIsEscaped, w, r)
					case "POST":
						s.handleStoreConfigRequest([0]string{}, elemIsEscaped, w, r)
					case "PUT":
						s.handleSetConfigRequest([0]string{}, elemIsEscaped, w, r)
					default:
						s.notAllowed(w, r, "GET,POST,PUT")
					}

					return
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'j': // Prefix: "jobs"
						if l := len("jobs"); len(elem) >= l && elem[0:l] == "jobs" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleGetJobsConfigRequest([0]string{}, elemIsEscaped, w, r)
							case "PUT":
								s.handleSetJobsConfigRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET,PUT")
							}

							return
						}
					case 'v': // Prefix: "views"
						if l := len("views"); len(elem) >= l && elem[0:l] == "views" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "DELETE":
								s.handleDeleteViewRequest([0]string{}, elemIsEscaped, w, r)
							case "GET":
								s.handleGetViewsRequest([0]string{}, elemIsEscaped, w, r)
							case "POST":
								s.handleAddViewRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "DELETE,GET,POST")
							}

							return
						}
					}
				}
			case 'i': // Prefix: "image/"
				if l := len("image/"); len(elem) >= l && elem[0:l] == "image/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "table"
				// Match until "/"
				idx := strings.IndexByte(elem, '/')
				if idx < 0 {
					idx = len(elem)
				}
				args[0] = elem[:idx]
				elem = elem[idx:]

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "field"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[1] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "search"
						// Leaf parameter
						args[2] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleGetImageRequest([3]string{
									args[0],
									args[1],
									args[2],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}
					}
				}
			case 'l': // Prefix: "log"
				if l := len("log"); len(elem) >= l && elem[0:l] == "log" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'i': // Prefix: "in"
					if l := len("in"); len(elem) >= l && elem[0:l] == "in" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetLoginSessionRequest([0]string{}, elemIsEscaped, w, r)
						case "POST":
							s.handlePushLoginSessionRequest([0]string{}, elemIsEscaped, w, r)
						case "PUT":
							s.handleLoginSessionRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET,POST,PUT")
						}

						return
					}
				case 'o': // Prefix: "off"
					if l := len("off"); len(elem) >= l && elem[0:l] == "off" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleRemoveSessionCompatRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}
				}
			case 'r': // Prefix: "rest/"
				if l := len("rest/"); len(elem) >= l && elem[0:l] == "rest/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'd': // Prefix: "database"
					if l := len("database"); len(elem) >= l && elem[0:l] == "database" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch r.Method {
						case "GET":
							s.handleGetDatabasesRequest([0]string{}, elemIsEscaped, w, r)
						case "POST":
							s.handlePostDatabaseRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET,POST")
						}

						return
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "table"
						// Match until "/"
						idx := strings.IndexByte(elem, '/')
						if idx < 0 {
							idx = len(elem)
						}
						args[0] = elem[:idx]
						elem = elem[idx:]

						if len(elem) == 0 {
							switch r.Method {
							case "DELETE":
								s.handleDeleteDatabaseRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							case "GET":
								s.handleDatabaseOperationRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							case "POST":
								s.handleDatabasePostOperationsRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							case "PUT":
								s.handlePutDatabaseResourceRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "DELETE,GET,POST,PUT")
							}

							return
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'c': // Prefix: "connection"
								if l := len("connection"); len(elem) >= l && elem[0:l] == "connection" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "DELETE":
										s.handleDisconnectTCPRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									case "GET":
										s.handleGetConnectionsRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "DELETE,GET")
									}

									return
								}
							case 'p': // Prefix: "permission"
								if l := len("permission"); len(elem) >= l && elem[0:l] == "permission" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch r.Method {
									case "DELETE":
										s.handleRemovePermissionRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									case "GET":
										s.handleGetPermissionRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									case "PUT":
										s.handleAdaptPermissionRequest([1]string{
											args[0],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "DELETE,GET,PUT")
									}

									return
								}
								switch elem[0] {
								case '/': // Prefix: "/"
									if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
										elem = elem[l:]
									} else {
										break
									}

									// Param: "resource"
									// Match until "/"
									idx := strings.IndexByte(elem, '/')
									if idx < 0 {
										idx = len(elem)
									}
									args[1] = elem[:idx]
									elem = elem[idx:]

									if len(elem) == 0 {
										switch r.Method {
										case "GET":
											s.handleListRBACResourceRequest([2]string{
												args[0],
												args[1],
											}, elemIsEscaped, w, r)
										default:
											s.notAllowed(w, r, "GET")
										}

										return
									}
									switch elem[0] {
									case '/': // Prefix: "/"
										if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
											elem = elem[l:]
										} else {
											break
										}

										// Param: "name"
										// Leaf parameter
										args[2] = elem
										elem = ""

										if len(elem) == 0 {
											// Leaf node.
											switch r.Method {
											case "DELETE":
												s.handleDeleteRBACResourceRequest([3]string{
													args[0],
													args[1],
													args[2],
												}, elemIsEscaped, w, r)
											case "PUT":
												s.handleAddRBACResourceRequest([3]string{
													args[0],
													args[1],
													args[2],
												}, elemIsEscaped, w, r)
											default:
												s.notAllowed(w, r, "DELETE,PUT")
											}

											return
										}
									}
								}
							case 's': // Prefix: "s"
								if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									break
								}
								switch elem[0] {
								case 'e': // Prefix: "essions"
									if l := len("essions"); len(elem) >= l && elem[0:l] == "essions" {
										elem = elem[l:]
									} else {
										break
									}

									if len(elem) == 0 {
										// Leaf node.
										switch r.Method {
										case "GET":
											s.handleGetDatabaseSessionsRequest([1]string{
												args[0],
											}, elemIsEscaped, w, r)
										default:
											s.notAllowed(w, r, "GET")
										}

										return
									}
								case 't': // Prefix: "tats"
									if l := len("tats"); len(elem) >= l && elem[0:l] == "tats" {
										elem = elem[l:]
									} else {
										break
									}

									if len(elem) == 0 {
										// Leaf node.
										switch r.Method {
										case "GET":
											s.handleGetDatabaseStatsRequest([1]string{
												args[0],
											}, elemIsEscaped, w, r)
										default:
											s.notAllowed(w, r, "GET")
										}

										return
									}
								}
							}
						}
					}
				case 'e': // Prefix: "env"
					if l := len("env"); len(elem) >= l && elem[0:l] == "env" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetEnvironmentsRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}
				case 'f': // Prefix: "file/"
					if l := len("file/"); len(elem) >= l && elem[0:l] == "file/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'b': // Prefix: "browse"
						if l := len("browse"); len(elem) >= l && elem[0:l] == "browse" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch r.Method {
							case "GET":
								s.handleBrowseListRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "path"
							// Leaf parameter
							args[0] = elem
							elem = ""

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleBrowseLocationRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}
						}
					}
					// Param: "path"
					// Leaf parameter
					args[0] = elem
					elem = ""

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "DELETE":
							s.handleDeleteFileLocationRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						case "GET":
							s.handleDownloadFileRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						case "POST":
							s.handleUploadFileRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						case "PUT":
							s.handleCreateDirectoryRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "DELETE,GET,POST,PUT")
						}

						return
					}
				case 'm': // Prefix: "m"
					if l := len("m"); len(elem) >= l && elem[0:l] == "m" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'a': // Prefix: "ap"
						if l := len("ap"); len(elem) >= l && elem[0:l] == "ap" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch r.Method {
							case "GET":
								s.handleListModellingRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "path"
							// Leaf parameter
							args[0] = elem
							elem = ""

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleSearchModellingRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}
						}
					case 'e': // Prefix: "etadata/view/"
						if l := len("etadata/view/"); len(elem) >= l && elem[0:l] == "etadata/view/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "table"
						// Leaf parameter
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							// Leaf node.
							switch r.Method {
							case "GET":
								s.handleGetMapMetadataRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}
					}
				case 's': // Prefix: "shutdown/"
					if l := len("shutdown/"); len(elem) >= l && elem[0:l] == "shutdown/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "hash"
					// Leaf parameter
					args[0] = elem
					elem = ""

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "PUT":
							s.handleShutdownServerRequest([1]string{
								args[0],
							}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "PUT")
						}

						return
					}
				case 't': // Prefix: "ta"
					if l := len("ta"); len(elem) >= l && elem[0:l] == "ta" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'b': // Prefix: "bles"
						if l := len("bles"); len(elem) >= l && elem[0:l] == "bles" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch r.Method {
							case "GET":
								s.handleListTablesRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET")
							}

							return
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "table"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[0] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									break
								}
								switch elem[0] {
								case 'f': // Prefix: "fields"
									if l := len("fields"); len(elem) >= l && elem[0:l] == "fields" {
										elem = elem[l:]
									} else {
										break
									}

									if len(elem) == 0 {
										// Leaf node.
										switch r.Method {
										case "GET":
											s.handleGetFieldsRequest([1]string{
												args[0],
											}, elemIsEscaped, w, r)
										default:
											s.notAllowed(w, r, "GET")
										}

										return
									}
								}
								// Param: "fields"
								// Match until "/"
								idx := strings.IndexByte(elem, '/')
								if idx < 0 {
									idx = len(elem)
								}
								args[1] = elem[:idx]
								elem = elem[idx:]

								if len(elem) == 0 {
									break
								}
								switch elem[0] {
								case '/': // Prefix: "/"
									if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
										elem = elem[l:]
									} else {
										break
									}

									// Param: "search"
									// Leaf parameter
									args[2] = elem
									elem = ""

									if len(elem) == 0 {
										// Leaf node.
										switch r.Method {
										case "GET":
											s.handleSearchTableRequest([3]string{
												args[0],
												args[1],
												args[2],
											}, elemIsEscaped, w, r)
										default:
											s.notAllowed(w, r, "GET")
										}

										return
									}
								}
							}
						}
					case 's': // Prefix: "sks"
						if l := len("sks"); len(elem) >= l && elem[0:l] == "sks" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch r.Method {
							case "GET":
								s.handleGetJobsRequest([0]string{}, elemIsEscaped, w, r)
							case "POST":
								s.handlePostJobRequest([0]string{}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "GET,POST")
							}

							return
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'r': // Prefix: "results"
								if l := len("results"); len(elem) >= l && elem[0:l] == "results" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetJobExecutionResultRequest([0]string{}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}
							}
							// Param: "jobName"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[0] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								switch r.Method {
								case "GET":
									s.handleGetJobFullInfoRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								case "PUT":
									s.handleTriggerJobRequest([1]string{
										args[0],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET,PUT")
								}

								return
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "jobId"
								// Leaf parameter
								args[1] = elem
								elem = ""

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "DELETE":
										s.handleDeleteJobResultRequest([2]string{
											args[0],
											args[1],
										}, elemIsEscaped, w, r)
									case "GET":
										s.handleGetJobResultRequest([2]string{
											args[0],
											args[1],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "DELETE,GET")
									}

									return
								}
							}
						}
					}
				case 'v': // Prefix: "view"
					if l := len("view"); len(elem) >= l && elem[0:l] == "view" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch r.Method {
						case "GET":
							s.handleGetMapsRequest([0]string{}, elemIsEscaped, w, r)
						case "POST":
							s.handleInsertMapFileRecordsRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET,POST")
						}

						return
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "table"
						// Match until "/"
						idx := strings.IndexByte(elem, '/')
						if idx < 0 {
							idx = len(elem)
						}
						args[0] = elem[:idx]
						elem = elem[idx:]

						if len(elem) == 0 {
							switch r.Method {
							case "POST":
								s.handleInsertRecordRequest([1]string{
									args[0],
								}, elemIsEscaped, w, r)
							default:
								s.notAllowed(w, r, "POST")
							}

							return
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "search"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[1] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								switch r.Method {
								case "DELETE":
									s.handleDeleteRecordsSearchedRequest([2]string{
										args[0],
										args[1],
									}, elemIsEscaped, w, r)
								case "GET":
									s.handleSearchRecordsFieldsRequest([2]string{
										args[0],
										args[1],
									}, elemIsEscaped, w, r)
								case "PUT":
									s.handleUpdateRecordsByFieldsRequest([2]string{
										args[0],
										args[1],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "DELETE,GET,PUT")
								}

								return
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "search"
								// Leaf parameter
								args[2] = elem
								elem = ""

								if len(elem) == 0 {
									// Leaf node.
									switch r.Method {
									case "GET":
										s.handleGetMapRecordsFieldsRequest([3]string{
											args[0],
											args[1],
											args[2],
										}, elemIsEscaped, w, r)
									default:
										s.notAllowed(w, r, "GET")
									}

									return
								}
							}
						}
					}
				}
			case 'v': // Prefix: "v"
				if l := len("v"); len(elem) >= l && elem[0:l] == "v" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'e': // Prefix: "ersion"
					if l := len("ersion"); len(elem) >= l && elem[0:l] == "ersion" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						// Leaf node.
						switch r.Method {
						case "GET":
							s.handleGetVersionRequest([0]string{}, elemIsEscaped, w, r)
						default:
							s.notAllowed(w, r, "GET")
						}

						return
					}
				case 'i': // Prefix: "ideo/"
					if l := len("ideo/"); len(elem) >= l && elem[0:l] == "ideo/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "table"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[0] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "field"
						// Match until "/"
						idx := strings.IndexByte(elem, '/')
						if idx < 0 {
							idx = len(elem)
						}
						args[1] = elem[:idx]
						elem = elem[idx:]

						if len(elem) == 0 {
							break
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "search"
							// Leaf parameter
							args[2] = elem
							elem = ""

							if len(elem) == 0 {
								// Leaf node.
								switch r.Method {
								case "GET":
									s.handleGetVideoRequest([3]string{
										args[0],
										args[1],
										args[2],
									}, elemIsEscaped, w, r)
								default:
									s.notAllowed(w, r, "GET")
								}

								return
							}
						}
					}
				}
			}
		}
	}
	log.Log.Debugf("Elem not found %s", elem)
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	operationID string
	pathPattern string
	count       int
	args        [3]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// PathPattern returns OpenAPI path.
func (r Route) PathPattern() string {
	return r.pathPattern
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
//
// Note: this method does not unescape path or handle reserved characters in path properly. Use FindPath instead.
func (s *Server) FindRoute(method, path string) (Route, bool) {
	return s.FindPath(method, &url.URL{Path: path})
}

// FindPath finds Route for given method and URL.
func (s *Server) FindPath(method string, u *url.URL) (r Route, _ bool) {
	var (
		elem = u.Path
		args = r.args
	)
	if rawPath := u.RawPath; rawPath != "" {
		if normalized, ok := uri.NormalizeEscapedPath(rawPath); ok {
			elem = normalized
		}
		defer func() {
			for i, arg := range r.args[:r.count] {
				if unescaped, err := url.PathUnescape(arg); err == nil {
					r.args[i] = unescaped
				}
			}
		}()
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/"
			if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				break
			}
			switch elem[0] {
			case 'a': // Prefix: "admin/access/"
				if l := len("admin/access/"); len(elem) >= l && elem[0:l] == "admin/access/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "role"
				// Leaf parameter
				args[0] = elem
				elem = ""

				if len(elem) == 0 {
					switch method {
					case "DELETE":
						// Leaf: DelAccess
						r.name = "DelAccess"
						r.operationID = "delAccess"
						r.pathPattern = "/admin/access/{role}"
						r.args = args
						r.count = 1
						return r, true
					case "GET":
						// Leaf: Access
						r.name = "Access"
						r.operationID = "access"
						r.pathPattern = "/admin/access/{role}"
						r.args = args
						r.count = 1
						return r, true
					case "POST":
						// Leaf: AddAccess
						r.name = "AddAccess"
						r.operationID = "addAccess"
						r.pathPattern = "/admin/access/{role}"
						r.args = args
						r.count = 1
						return r, true
					default:
						return
					}
				}
			case 'b': // Prefix: "binary/"
				if l := len("binary/"); len(elem) >= l && elem[0:l] == "binary/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "table"
				// Match until "/"
				idx := strings.IndexByte(elem, '/')
				if idx < 0 {
					idx = len(elem)
				}
				args[0] = elem[:idx]
				elem = elem[idx:]

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "field"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[1] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "search"
						// Leaf parameter
						args[2] = elem
						elem = ""

						if len(elem) == 0 {
							switch method {
							case "GET":
								// Leaf: GetLobByMap
								r.name = "GetLobByMap"
								r.operationID = "getLobByMap"
								r.pathPattern = "/binary/{table}/{field}/{search}"
								r.args = args
								r.count = 3
								return r, true
							case "PUT":
								// Leaf: UpdateLobByMap
								r.name = "UpdateLobByMap"
								r.operationID = "updateLobByMap"
								r.pathPattern = "/binary/{table}/{field}/{search}"
								r.args = args
								r.count = 3
								return r, true
							default:
								return
							}
						}
					}
				}
			case 'c': // Prefix: "config"
				if l := len("config"); len(elem) >= l && elem[0:l] == "config" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					switch method {
					case "GET":
						r.name = "GetConfig"
						r.operationID = "getConfig"
						r.pathPattern = "/config"
						r.args = args
						r.count = 0
						return r, true
					case "POST":
						r.name = "StoreConfig"
						r.operationID = "storeConfig"
						r.pathPattern = "/config"
						r.args = args
						r.count = 0
						return r, true
					case "PUT":
						r.name = "SetConfig"
						r.operationID = "setConfig"
						r.pathPattern = "/config"
						r.args = args
						r.count = 0
						return r, true
					default:
						return
					}
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'j': // Prefix: "jobs"
						if l := len("jobs"); len(elem) >= l && elem[0:l] == "jobs" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "GET":
								// Leaf: GetJobsConfig
								r.name = "GetJobsConfig"
								r.operationID = "getJobsConfig"
								r.pathPattern = "/config/jobs"
								r.args = args
								r.count = 0
								return r, true
							case "PUT":
								// Leaf: SetJobsConfig
								r.name = "SetJobsConfig"
								r.operationID = "setJobsConfig"
								r.pathPattern = "/config/jobs"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
					case 'v': // Prefix: "views"
						if l := len("views"); len(elem) >= l && elem[0:l] == "views" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "DELETE":
								// Leaf: DeleteView
								r.name = "DeleteView"
								r.operationID = "deleteView"
								r.pathPattern = "/config/views"
								r.args = args
								r.count = 0
								return r, true
							case "GET":
								// Leaf: GetViews
								r.name = "GetViews"
								r.operationID = "getViews"
								r.pathPattern = "/config/views"
								r.args = args
								r.count = 0
								return r, true
							case "POST":
								// Leaf: AddView
								r.name = "AddView"
								r.operationID = "addView"
								r.pathPattern = "/config/views"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
					}
				}
			case 'i': // Prefix: "image/"
				if l := len("image/"); len(elem) >= l && elem[0:l] == "image/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "table"
				// Match until "/"
				idx := strings.IndexByte(elem, '/')
				if idx < 0 {
					idx = len(elem)
				}
				args[0] = elem[:idx]
				elem = elem[idx:]

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case '/': // Prefix: "/"
					if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "field"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[1] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "search"
						// Leaf parameter
						args[2] = elem
						elem = ""

						if len(elem) == 0 {
							switch method {
							case "GET":
								// Leaf: GetImage
								r.name = "GetImage"
								r.operationID = "getImage"
								r.pathPattern = "/image/{table}/{field}/{search}"
								r.args = args
								r.count = 3
								return r, true
							default:
								return
							}
						}
					}
				}
			case 'l': // Prefix: "log"
				if l := len("log"); len(elem) >= l && elem[0:l] == "log" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'i': // Prefix: "in"
					if l := len("in"); len(elem) >= l && elem[0:l] == "in" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							// Leaf: GetLoginSession
							r.name = "GetLoginSession"
							r.operationID = "getLoginSession"
							r.pathPattern = "/login"
							r.args = args
							r.count = 0
							return r, true
						case "POST":
							// Leaf: PushLoginSession
							r.name = "PushLoginSession"
							r.operationID = "pushLoginSession"
							r.pathPattern = "/login"
							r.args = args
							r.count = 0
							return r, true
						case "PUT":
							// Leaf: LoginSession
							r.name = "LoginSession"
							r.operationID = "loginSession"
							r.pathPattern = "/login"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				case 'o': // Prefix: "off"
					if l := len("off"); len(elem) >= l && elem[0:l] == "off" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							// Leaf: RemoveSessionCompat
							r.name = "RemoveSessionCompat"
							r.operationID = "removeSessionCompat"
							r.pathPattern = "/logoff"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				}
			case 'r': // Prefix: "rest/"
				if l := len("rest/"); len(elem) >= l && elem[0:l] == "rest/" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'd': // Prefix: "database"
					if l := len("database"); len(elem) >= l && elem[0:l] == "database" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							r.name = "GetDatabases"
							r.operationID = "getDatabases"
							r.pathPattern = "/rest/database"
							r.args = args
							r.count = 0
							return r, true
						case "POST":
							r.name = "PostDatabase"
							r.operationID = "postDatabase"
							r.pathPattern = "/rest/database"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "table"
						// Match until "/"
						idx := strings.IndexByte(elem, '/')
						if idx < 0 {
							idx = len(elem)
						}
						args[0] = elem[:idx]
						elem = elem[idx:]

						if len(elem) == 0 {
							switch method {
							case "DELETE":
								r.name = "DeleteDatabase"
								r.operationID = "deleteDatabase"
								r.pathPattern = "/rest/database/{table_operation}"
								r.args = args
								r.count = 1
								return r, true
							case "GET":
								r.name = "DatabaseOperation"
								r.operationID = "databaseOperation"
								r.pathPattern = "/rest/database/{table_operation}"
								r.args = args
								r.count = 1
								return r, true
							case "POST":
								r.name = "DatabasePostOperations"
								r.operationID = "databasePostOperations"
								r.pathPattern = "/rest/database/{table_operation}"
								r.args = args
								r.count = 1
								return r, true
							case "PUT":
								r.name = "PutDatabaseResource"
								r.operationID = "putDatabaseResource"
								r.pathPattern = "/rest/database/{table_operation}"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'c': // Prefix: "connection"
								if l := len("connection"); len(elem) >= l && elem[0:l] == "connection" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch method {
									case "DELETE":
										// Leaf: DisconnectTCP
										r.name = "DisconnectTCP"
										r.operationID = "disconnectTCP"
										r.pathPattern = "/rest/database/{table}/connection"
										r.args = args
										r.count = 1
										return r, true
									case "GET":
										// Leaf: GetConnections
										r.name = "GetConnections"
										r.operationID = "getConnections"
										r.pathPattern = "/rest/database/{table}/connection"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}
							case 'p': // Prefix: "permission"
								if l := len("permission"); len(elem) >= l && elem[0:l] == "permission" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch method {
									case "DELETE":
										r.name = "RemovePermission"
										r.operationID = "removePermission"
										r.pathPattern = "/rest/database/{table}/permission"
										r.args = args
										r.count = 1
										return r, true
									case "GET":
										r.name = "GetPermission"
										r.operationID = "getPermission"
										r.pathPattern = "/rest/database/{table}/permission"
										r.args = args
										r.count = 1
										return r, true
									case "PUT":
										r.name = "AdaptPermission"
										r.operationID = "adaptPermission"
										r.pathPattern = "/rest/database/{table}/permission"
										r.args = args
										r.count = 1
										return r, true
									default:
										return
									}
								}
								switch elem[0] {
								case '/': // Prefix: "/"
									if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
										elem = elem[l:]
									} else {
										break
									}

									// Param: "resource"
									// Match until "/"
									idx := strings.IndexByte(elem, '/')
									if idx < 0 {
										idx = len(elem)
									}
									args[1] = elem[:idx]
									elem = elem[idx:]

									if len(elem) == 0 {
										switch method {
										case "GET":
											r.name = "ListRBACResource"
											r.operationID = "listRBACResource"
											r.pathPattern = "/rest/database/{table}/permission/{resource}"
											r.args = args
											r.count = 2
											return r, true
										default:
											return
										}
									}
									switch elem[0] {
									case '/': // Prefix: "/"
										if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
											elem = elem[l:]
										} else {
											break
										}

										// Param: "name"
										// Leaf parameter
										args[2] = elem
										elem = ""

										if len(elem) == 0 {
											switch method {
											case "DELETE":
												// Leaf: DeleteRBACResource
												r.name = "DeleteRBACResource"
												r.operationID = "deleteRBACResource"
												r.pathPattern = "/rest/database/{table}/permission/{resource}/{name}"
												r.args = args
												r.count = 3
												return r, true
											case "PUT":
												// Leaf: AddRBACResource
												r.name = "AddRBACResource"
												r.operationID = "addRBACResource"
												r.pathPattern = "/rest/database/{table}/permission/{resource}/{name}"
												r.args = args
												r.count = 3
												return r, true
											default:
												return
											}
										}
									}
								}
							case 's': // Prefix: "s"
								if l := len("s"); len(elem) >= l && elem[0:l] == "s" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									break
								}
								switch elem[0] {
								case 'e': // Prefix: "essions"
									if l := len("essions"); len(elem) >= l && elem[0:l] == "essions" {
										elem = elem[l:]
									} else {
										break
									}

									if len(elem) == 0 {
										switch method {
										case "GET":
											// Leaf: GetDatabaseSessions
											r.name = "GetDatabaseSessions"
											r.operationID = "getDatabaseSessions"
											r.pathPattern = "/rest/database/{table}/sessions"
											r.args = args
											r.count = 1
											return r, true
										default:
											return
										}
									}
								case 't': // Prefix: "tats"
									if l := len("tats"); len(elem) >= l && elem[0:l] == "tats" {
										elem = elem[l:]
									} else {
										break
									}

									if len(elem) == 0 {
										switch method {
										case "GET":
											// Leaf: GetDatabaseStats
											r.name = "GetDatabaseStats"
											r.operationID = "getDatabaseStats"
											r.pathPattern = "/rest/database/{table}/stats"
											r.args = args
											r.count = 1
											return r, true
										default:
											return
										}
									}
								}
							}
						}
					}
				case 'e': // Prefix: "env"
					if l := len("env"); len(elem) >= l && elem[0:l] == "env" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							// Leaf: GetEnvironments
							r.name = "GetEnvironments"
							r.operationID = "getEnvironments"
							r.pathPattern = "/rest/env"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				case 'f': // Prefix: "file/"
					if l := len("file/"); len(elem) >= l && elem[0:l] == "file/" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'b': // Prefix: "browse"
						if l := len("browse"); len(elem) >= l && elem[0:l] == "browse" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "GET":
								r.name = "BrowseList"
								r.operationID = "browseList"
								r.pathPattern = "/rest/file/browse"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "path"
							// Leaf parameter
							args[0] = elem
							elem = ""

							if len(elem) == 0 {
								switch method {
								case "GET":
									// Leaf: BrowseLocation
									r.name = "BrowseLocation"
									r.operationID = "browseLocation"
									r.pathPattern = "/rest/file/browse/{path}"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}
						}
					}
					// Param: "path"
					// Leaf parameter
					args[0] = elem
					elem = ""

					if len(elem) == 0 {
						switch method {
						case "DELETE":
							// Leaf: DeleteFileLocation
							r.name = "DeleteFileLocation"
							r.operationID = "deleteFileLocation"
							r.pathPattern = "/rest/file/{path}"
							r.args = args
							r.count = 1
							return r, true
						case "GET":
							// Leaf: DownloadFile
							r.name = "DownloadFile"
							r.operationID = "downloadFile"
							r.pathPattern = "/rest/file/{path}"
							r.args = args
							r.count = 1
							return r, true
						case "POST":
							// Leaf: UploadFile
							r.name = "UploadFile"
							r.operationID = "uploadFile"
							r.pathPattern = "/rest/file/{path}"
							r.args = args
							r.count = 1
							return r, true
						case "PUT":
							// Leaf: CreateDirectory
							r.name = "CreateDirectory"
							r.operationID = "createDirectory"
							r.pathPattern = "/rest/file/{path}"
							r.args = args
							r.count = 1
							return r, true
						default:
							return
						}
					}
				case 'm': // Prefix: "m"
					if l := len("m"); len(elem) >= l && elem[0:l] == "m" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'a': // Prefix: "ap"
						if l := len("ap"); len(elem) >= l && elem[0:l] == "ap" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "GET":
								r.name = "ListModelling"
								r.operationID = "listModelling"
								r.pathPattern = "/rest/map"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "path"
							// Leaf parameter
							args[0] = elem
							elem = ""

							if len(elem) == 0 {
								switch method {
								case "GET":
									// Leaf: SearchModelling
									r.name = "SearchModelling"
									r.operationID = "searchModelling"
									r.pathPattern = "/rest/map/{path}"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}
						}
					case 'e': // Prefix: "etadata/view/"
						if l := len("etadata/view/"); len(elem) >= l && elem[0:l] == "etadata/view/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "table"
						// Leaf parameter
						args[0] = elem
						elem = ""

						if len(elem) == 0 {
							switch method {
							case "GET":
								// Leaf: GetMapMetadata
								r.name = "GetMapMetadata"
								r.operationID = "getMapMetadata"
								r.pathPattern = "/rest/metadata/view/{table}"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}
					}
				case 's': // Prefix: "shutdown/"
					if l := len("shutdown/"); len(elem) >= l && elem[0:l] == "shutdown/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "hash"
					// Leaf parameter
					args[0] = elem
					elem = ""

					if len(elem) == 0 {
						switch method {
						case "PUT":
							// Leaf: ShutdownServer
							r.name = "ShutdownServer"
							r.operationID = "shutdownServer"
							r.pathPattern = "/rest/shutdown/{hash}"
							r.args = args
							r.count = 1
							return r, true
						default:
							return
						}
					}
				case 't': // Prefix: "ta"
					if l := len("ta"); len(elem) >= l && elem[0:l] == "ta" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case 'b': // Prefix: "bles"
						if l := len("bles"); len(elem) >= l && elem[0:l] == "bles" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "GET":
								r.name = "ListTables"
								r.operationID = "listTables"
								r.pathPattern = "/rest/tables"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "table"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[0] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									break
								}
								switch elem[0] {
								case 'f': // Prefix: "fields"
									if l := len("fields"); len(elem) >= l && elem[0:l] == "fields" {
										elem = elem[l:]
									} else {
										break
									}

									if len(elem) == 0 {
										switch method {
										case "GET":
											// Leaf: GetFields
											r.name = "GetFields"
											r.operationID = "getFields"
											r.pathPattern = "/rest/tables/{table}/fields"
											r.args = args
											r.count = 1
											return r, true
										default:
											return
										}
									}
								}
								// Param: "fields"
								// Match until "/"
								idx := strings.IndexByte(elem, '/')
								if idx < 0 {
									idx = len(elem)
								}
								args[1] = elem[:idx]
								elem = elem[idx:]

								if len(elem) == 0 {
									break
								}
								switch elem[0] {
								case '/': // Prefix: "/"
									if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
										elem = elem[l:]
									} else {
										break
									}

									// Param: "search"
									// Leaf parameter
									args[2] = elem
									elem = ""

									if len(elem) == 0 {
										switch method {
										case "GET":
											// Leaf: SearchTable
											r.name = "SearchTable"
											r.operationID = "searchTable"
											r.pathPattern = "/rest/tables/{table}/{fields}/{search}"
											r.args = args
											r.count = 3
											return r, true
										default:
											return
										}
									}
								}
							}
						}
					case 's': // Prefix: "sks"
						if l := len("sks"); len(elem) >= l && elem[0:l] == "sks" {
							elem = elem[l:]
						} else {
							break
						}

						if len(elem) == 0 {
							switch method {
							case "GET":
								r.name = "GetJobs"
								r.operationID = "getJobs"
								r.pathPattern = "/rest/tasks"
								r.args = args
								r.count = 0
								return r, true
							case "POST":
								r.name = "PostJob"
								r.operationID = "postJob"
								r.pathPattern = "/rest/tasks"
								r.args = args
								r.count = 0
								return r, true
							default:
								return
							}
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							if len(elem) == 0 {
								break
							}
							switch elem[0] {
							case 'r': // Prefix: "results"
								if l := len("results"); len(elem) >= l && elem[0:l] == "results" {
									elem = elem[l:]
								} else {
									break
								}

								if len(elem) == 0 {
									switch method {
									case "GET":
										// Leaf: GetJobExecutionResult
										r.name = "GetJobExecutionResult"
										r.operationID = "getJobExecutionResult"
										r.pathPattern = "/rest/tasks/results"
										r.args = args
										r.count = 0
										return r, true
									default:
										return
									}
								}
							}
							// Param: "jobName"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[0] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								switch method {
								case "GET":
									r.name = "GetJobFullInfo"
									r.operationID = "getJobFullInfo"
									r.pathPattern = "/rest/tasks/{jobName}"
									r.args = args
									r.count = 1
									return r, true
								case "PUT":
									r.name = "TriggerJob"
									r.operationID = "triggerJob"
									r.pathPattern = "/rest/tasks/{jobName}"
									r.args = args
									r.count = 1
									return r, true
								default:
									return
								}
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "jobId"
								// Leaf parameter
								args[1] = elem
								elem = ""

								if len(elem) == 0 {
									switch method {
									case "DELETE":
										// Leaf: DeleteJobResult
										r.name = "DeleteJobResult"
										r.operationID = "deleteJobResult"
										r.pathPattern = "/rest/tasks/{jobName}/{jobId}"
										r.args = args
										r.count = 2
										return r, true
									case "GET":
										// Leaf: GetJobResult
										r.name = "GetJobResult"
										r.operationID = "getJobResult"
										r.pathPattern = "/rest/tasks/{jobName}/{jobId}"
										r.args = args
										r.count = 2
										return r, true
									default:
										return
									}
								}
							}
						}
					}
				case 'v': // Prefix: "view"
					if l := len("view"); len(elem) >= l && elem[0:l] == "view" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							r.name = "GetMaps"
							r.operationID = "getMaps"
							r.pathPattern = "/rest/view"
							r.args = args
							r.count = 0
							return r, true
						case "POST":
							r.name = "InsertMapFileRecords"
							r.operationID = "insertMapFileRecords"
							r.pathPattern = "/rest/view"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "table"
						// Match until "/"
						idx := strings.IndexByte(elem, '/')
						if idx < 0 {
							idx = len(elem)
						}
						args[0] = elem[:idx]
						elem = elem[idx:]

						if len(elem) == 0 {
							switch method {
							case "POST":
								r.name = "InsertRecord"
								r.operationID = "insertRecord"
								r.pathPattern = "/rest/view/{table}"
								r.args = args
								r.count = 1
								return r, true
							default:
								return
							}
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "search"
							// Match until "/"
							idx := strings.IndexByte(elem, '/')
							if idx < 0 {
								idx = len(elem)
							}
							args[1] = elem[:idx]
							elem = elem[idx:]

							if len(elem) == 0 {
								switch method {
								case "DELETE":
									r.name = "DeleteRecordsSearched"
									r.operationID = "deleteRecordsSearched"
									r.pathPattern = "/rest/view/{table}/{search}"
									r.args = args
									r.count = 2
									return r, true
								case "GET":
									r.name = "SearchRecordsFields"
									r.operationID = "searchRecordsFields"
									r.pathPattern = "/rest/view/{table}/{search}"
									r.args = args
									r.count = 2
									return r, true
								case "PUT":
									r.name = "UpdateRecordsByFields"
									r.operationID = "updateRecordsByFields"
									r.pathPattern = "/rest/view/{table}/{search}"
									r.args = args
									r.count = 2
									return r, true
								default:
									return
								}
							}
							switch elem[0] {
							case '/': // Prefix: "/"
								if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
									elem = elem[l:]
								} else {
									break
								}

								// Param: "search"
								// Leaf parameter
								args[2] = elem
								elem = ""

								if len(elem) == 0 {
									switch method {
									case "GET":
										// Leaf: GetMapRecordsFields
										r.name = "GetMapRecordsFields"
										r.operationID = "getMapRecordsFields"
										r.pathPattern = "/rest/view/{table}/{fields}/{search}"
										r.args = args
										r.count = 3
										return r, true
									default:
										return
									}
								}
							}
						}
					}
				}
			case 'v': // Prefix: "v"
				if l := len("v"); len(elem) >= l && elem[0:l] == "v" {
					elem = elem[l:]
				} else {
					break
				}

				if len(elem) == 0 {
					break
				}
				switch elem[0] {
				case 'e': // Prefix: "ersion"
					if l := len("ersion"); len(elem) >= l && elem[0:l] == "ersion" {
						elem = elem[l:]
					} else {
						break
					}

					if len(elem) == 0 {
						switch method {
						case "GET":
							// Leaf: GetVersion
							r.name = "GetVersion"
							r.operationID = "getVersion"
							r.pathPattern = "/version"
							r.args = args
							r.count = 0
							return r, true
						default:
							return
						}
					}
				case 'i': // Prefix: "ideo/"
					if l := len("ideo/"); len(elem) >= l && elem[0:l] == "ideo/" {
						elem = elem[l:]
					} else {
						break
					}

					// Param: "table"
					// Match until "/"
					idx := strings.IndexByte(elem, '/')
					if idx < 0 {
						idx = len(elem)
					}
					args[0] = elem[:idx]
					elem = elem[idx:]

					if len(elem) == 0 {
						break
					}
					switch elem[0] {
					case '/': // Prefix: "/"
						if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
							elem = elem[l:]
						} else {
							break
						}

						// Param: "field"
						// Match until "/"
						idx := strings.IndexByte(elem, '/')
						if idx < 0 {
							idx = len(elem)
						}
						args[1] = elem[:idx]
						elem = elem[idx:]

						if len(elem) == 0 {
							break
						}
						switch elem[0] {
						case '/': // Prefix: "/"
							if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
								elem = elem[l:]
							} else {
								break
							}

							// Param: "search"
							// Leaf parameter
							args[2] = elem
							elem = ""

							if len(elem) == 0 {
								switch method {
								case "GET":
									// Leaf: GetVideo
									r.name = "GetVideo"
									r.operationID = "getVideo"
									r.pathPattern = "/video/{table}/{field}/{search}"
									r.args = args
									r.count = 3
									return r, true
								default:
									return
								}
							}
						}
					}
				}
			}
		}
	}
	return r, false
}
