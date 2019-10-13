// This file is part of ezBastion.

//     ezBastion is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     ezBastion is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.

//     You should have received a copy of the GNU Affero General Public License
//     along with ezBastion.  If not, see <https://www.gnu.org/licenses/>.

package exec

// type EzbJob struct {
// 	Name      string
// 	Params    map[string]string
// 	Async     bool
// 	Path      string
// 	Xtrack    string
// 	Token     string
// 	Requester string
// 	Checksum  string
// }

type EzbJobs struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Enable   bool   `json:"enable"`
	Comment  string `json:"comment"`
	Checksum string `json:"checksum"`
	Path     string `json:"path"`
	Cache    int    `json:"cache"`
	Output   string `json:"output"`
}

type EzbParams struct {
	Data map[string]string `json:"data"`
	Meta EzbParamMeta      `json:"meta"`
}

type EzbParamMeta struct {
	Job EzbJobs `json:"job"`
}
